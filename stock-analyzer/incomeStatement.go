package main

import (
	"github.com/gocolly/colly"
)

type IncomeStatement struct {
	revenue                 string
	revenuePerShare         string
	grossProfit             string
	quarterlyEarningsGrowth string
}

func scrapeIncomeStatement(symbol string) *IncomeStatement {
	var incomeStatement IncomeStatement

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
			case "Revenue (ttm)":
				incomeStatement.revenue = dataSlice[1]
			case "Revenue Per Share (ttm)":
				incomeStatement.revenuePerShare = dataSlice[1]
			case "Gross Profit (ttm)":
				incomeStatement.grossProfit = dataSlice[1]
			case "Quarterly Earnings Growth (yoy)":
				incomeStatement.quarterlyEarningsGrowth = dataSlice[1]
			}
		})
	})

	c.Visit("https://ca.finance.yahoo.com/quote/" + symbol + "/key-statistics?p=" + symbol)
	c.Wait()
	return &incomeStatement
}
