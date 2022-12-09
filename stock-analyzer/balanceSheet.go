package main

import (
	"strconv"

	"github.com/gocolly/colly"
)

type BalanceSheet struct {
	totalCash         string
	totalDebt         string
	bookValuePerShare float64
}

func ScrapeBalanceSheet(symbol string) *BalanceSheet {
	var balanceSheet BalanceSheet
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

	c.OnHTML("tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			dataSlice := []string{}
			el.ForEach("td", func(_ int, el *colly.HTMLElement) {
				dataSlice = append(dataSlice, el.Text)
			})
			switch dataSlice[0] {
			case "Total Cash (mrq)":
				balanceSheet.totalCash = dataSlice[1]
			case "Total Debt (mrq)":
				balanceSheet.totalDebt = dataSlice[1] //strconv.ParseFloat(dataSlice[1], 8)
			case "Total Cash Per Share (mrq)":
				balanceSheet.bookValuePerShare, _ = strconv.ParseFloat(dataSlice[1], 64)
			}
		})
	})

	c.Visit("https://ca.finance.yahoo.com/quote/" + symbol + "/key-statistics?p=" + symbol)
	c.Wait()

	return &balanceSheet
}
