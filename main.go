// Package main x-stock is my stock bot
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/x-stock/exportor"
	"github.com/axiaoxin-com/x-stock/parser"
)

func main() {
	logging.SetLevel("info")
	ctx := context.Background()
	stocks, err := parser.AutoFilterStocks(ctx)
	if err != nil {
		logging.Fatal(ctx, err.Error())
	}
	data := exportor.New(ctx, stocks)
	data.SortByROE()
	filename := fmt.Sprintf("./docs/selected_stocks_%s.json", time.Now().Format("20060102"))
	c, err := data.ExportJSON(ctx, filename)
	if err != nil {
		logging.Fatal(ctx, err.Error())
	}
	logging.Debug(ctx, "json content:"+string(c))
	filename = fmt.Sprintf("./docs/selected_stocks_%s.csv", time.Now().Format("20060102"))
	c, err = data.ExportCSV(ctx, filename)
	if err != nil {
		logging.Fatal(ctx, err.Error())
	}
	logging.Debug(ctx, "csv content:"+string(c))
}
