package main

import "github.com/piquette/finance-go/quote"

type StockPriceHistory struct {
	yearChange          float64
	yearHigh            float64
	yearLow             float64
	fiftyDayMovingAvg   float64
	twoHundredMovingAvg float64
}

func scrapePriceHistory(symbol string) *StockPriceHistory {
	var stockPriceHistory StockPriceHistory
	stockHistoryPage, _ := quote.Get(symbol)

	stockPriceHistory.yearChange = stockHistoryPage.FiftyTwoWeekHighChangePercent
	stockPriceHistory.yearHigh = stockHistoryPage.FiftyTwoWeekHigh
	stockPriceHistory.yearLow = stockHistoryPage.FiftyTwoWeekLow
	stockPriceHistory.fiftyDayMovingAvg = stockHistoryPage.FiftyDayAverage
	stockPriceHistory.twoHundredMovingAvg = stockHistoryPage.TwoHundredDayAverage

	return &stockPriceHistory
}
