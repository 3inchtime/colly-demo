package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly"
)

type news struct {
	Title     string
	URL       string
	Contents  string
	CrawledAt time.Time
}

func main() {
	c := colly.NewCollector(
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36"),
	)

	c.Limit(&colly.LimitRule{
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

	detailCollector := c.Clone()

	detailCollector.OnResponse(func(r *colly.Response) {
		fmt.Println(r.StatusCode)
	})

	detailCollector.OnHTML("body", func(e *colly.HTMLElement) {
		n := news{}
		n.Title = e.ChildText("div[class=Mid2L_tit]>h1")
		n.URL = e.Request.URL.String()
		n.Contents = e.ChildText("div[class=Mid2L_con]>p")
		n.CrawledAt = time.Now()
		log.Println(n)
	})

	// Find and visit next page links
	c.OnHTML("a[class=tt]", func(e *colly.HTMLElement) {
		detailCollector.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit("https://www.gamersky.com/news/")

	c.Wait()
	detailCollector.Wait()
}
