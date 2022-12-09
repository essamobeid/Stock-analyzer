package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

type StockPriceHistory struct {
	yearChange          string
	yearHigh            int
	yearLow             int
	fiftyDayMovingAvg   int
	twoHundredMovingAvg int
}

type IncomeStatement struct {
	revenue                 int
	revenuePerShare         int
	grossProfit             int
	QuarterlyEarningsGrowth string
}

type BalanceSheet struct {
	totalCash         int
	totalDebt         int
	bookValuePerShare int
}

type CashFlow struct {
	operatingCashFlow   int
	LeveredFreeCashFlow int
}

func main() {
	flag.Parse()

	if len(flag.Args()) == 0 {
		logrus.Fatalf("Input one stock symbol", os.Args[0])
	}

	symbol := flag.Args()[0]
	summary := parseSummary(symbol)
	fmt.Println("Company Name:", summary.name)
	fmt.Println("Stock Price:", summary.ask)
	fmt.Println("Market Cap:", summary.marketCap)
	fmt.Println("PE Ration:", summary.peRation)
	fmt.Println("Daily Volume:", summary.volume)
	fmt.Println("Average Volume:", summary.avgVolume)
	fmt.Println("Earnings Date:", summary.earningDate)

}

// func printData(data [][]string, headers []string) {
// 	fmt.Println("Printing")
// 	fmt.Println(headers)
// 	for _, value := range data {
// 		fmt.Println(value)
// 	}
// }

func parseStock(symbol string) {
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
			log.Println(el.Text)

			if dataSlice[0] == "Previous Close" {
				fmt.Printf("Previous Price: $%v\n", dataSlice[1])
			}
		})
	})

	c.Visit("https://ca.finance.yahoo.com/quote/" + symbol)
	c.Wait()
}

func parseStatistics(symbol string) {
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
			log.Println(el.Text)

			if dataSlice[0] == "Enterprise Value " {
				log.Println("----------")
				fmt.Printf("Market Cap (intraday): $%v\n", dataSlice[1])
			}
		})
	})

	c.Visit("https://ca.finance.yahoo.com/quote/AAPL/key-statistics?p=AAPL")
	c.Wait()
}
