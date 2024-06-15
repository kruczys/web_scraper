package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
	"strconv"
	"strings"
)

type quote struct {
	text   string
	author string
	tags   []string
}

func scrapeSinglePage(url string, c *colly.Collector) []quote {
	var quotes []quote

	c.OnHTML("div.col-md-8", func(e *colly.HTMLElement) {
		e.ForEach("div.quote", func(i int, element *colly.HTMLElement) {
			var quote quote
			element.ForEach("span.text", func(_ int, quoteText *colly.HTMLElement) {
				quote.text = quoteText.Text
			})
			element.ForEach("small.author", func(_ int, author *colly.HTMLElement) {
				quote.author = author.Text
			})
			element.ForEach("div.tags", func(_ int, tagContainer *colly.HTMLElement) {
				tagContainer.ForEach("a.tag", func(_ int, tag *colly.HTMLElement) {
					quote.tags = append(quote.tags, tag.Text)
				})
			})
			quotes = append(quotes, quote)
		})
	})

	if err := c.Visit(url); err != nil {
		log.Fatal(err)
	}

	return quotes
}

func scrapeAllPages(url string, c *colly.Collector, firstPage, lastPage int) []quote {
	var quotes []quote
	for i := firstPage; i <= lastPage; i++ {
		url := fmt.Sprintf("%s/%d/", url, i)
		quotes = append(quotes, scrapeSinglePage(url, c)...)
	}
	return quotes
}

func main() {
	c := colly.NewCollector()
	result := scrapeAllPages("https://quotes.toscrape.com/page/", c, 1, 10)

	file, err := os.Create("quotes.csv")

	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}(file)

	if err != nil {
		log.Fatal(err)
	}

	csvWriter := csv.NewWriter(file)

	for index, quote := range result {
		csvRow := []string{strconv.Itoa(index), quote.text, quote.author, strings.Join(quote.tags, ", ")}
		if err := csvWriter.Write(csvRow); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("job finished")
}
