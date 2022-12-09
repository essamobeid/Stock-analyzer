package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gocolly/colly"
	"github.com/piquette/finance-go/quote"
)

type Summary struct {
	name        string
	ask         float64
	daysRange   string
	volume      float64
	avgVolume   float64
	marketCap   string
	peRation    float64
	earningDate string
}

func scrapeSummary(symbol string) *Summary {
	var summary Summary
	summaryPage, _ := quote.Get(symbol)
	summary.name = summaryPage.ShortName
	if summaryPage.Ask == 0 {
		summary.ask = summaryPage.RegularMarketPreviousClose
	} else {
		summary.ask = summaryPage.Ask
	}
	summary.ask = summaryPage.Ask
	summary.volume = float64(summaryPage.RegularMarketVolume)
	summary.avgVolume = float64(summaryPage.AverageDailyVolume3Month)

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.72 Safari/537.36"),
		colly.AllowedDomains("ca.finance.yahoo.com"),
		colly.MaxBodySize(0),
		colly.AllowURLRevisit(),
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		// Delay:       500 * time.Millisecond,
	})
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting ", r.URL.String())
	})

	c.OnHTML("tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			dataSlice := []string{}
			el.ForEach("td", func(_ int, el *colly.HTMLElement) {
				dataSlice = append(dataSlice, el.Text)
			})
			switch dataSlice[0] {
			case "Day's Range":
				summary.daysRange = dataSlice[1]
			case "Market Cap":
				summary.marketCap = dataSlice[1] //strconv.ParseFloat(dataSlice[1], 8)
			case "PE Ratio (TTM)":
				summary.peRation, _ = strconv.ParseFloat(dataSlice[1], 64)
			case "Earnings Date":
				summary.earningDate = dataSlice[1]
			}
			if dataSlice[0] == "Previous Close" {
				fmt.Printf("Previous Price: $%v\n", dataSlice[1])
			}
		})
	})

	c.Visit("https://ca.finance.yahoo.com/quote/" + symbol)
	c.Wait()

	return &summary
}
