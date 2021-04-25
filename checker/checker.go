// Package checker 对给定股票名/股票代码进行检测
package checker

import (
	"context"
	"fmt"
	"os"

	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/x-stock/core"
	"github.com/axiaoxin-com/x-stock/datacenter"
	"github.com/axiaoxin-com/x-stock/datacenter/eastmoney"
	"github.com/axiaoxin-com/x-stock/model"
	"github.com/olekukonko/tablewriter"
)

// Check 对给定名称或代码进行检测，输出检测结果
func Check(ctx context.Context, keyword string) {
	results, err := datacenter.QQ.KeywordSearch(ctx, keyword)
	if err != nil {
		logging.Fatal(ctx, err.Error())
	}
	if len(results) == 0 {
		logging.Fatal(ctx, "invalid keyword:"+keyword)
	}
	result := results[0]
	logging.Infof(ctx, "search result:%+v", result)
	filter := eastmoney.DefaultFilter
	filter.SpecialSecurityCode = result.SecurityCode
	stocks, err := datacenter.EastMoney.QuerySelectedStocksWithFilter(ctx, filter)
	if err != nil {
		logging.Fatal(ctx, "QuerySelectedStocksWithFilter error:"+err.Error())
	}
	if len(stocks) == 0 {
		logging.Fatal(ctx, "no stock data")
	}
	stock, err := model.NewStock(ctx, stocks[0], false)
	if err != nil {
		logging.Fatal(ctx, "NewStock error:"+err.Error())
	}
	checker := core.NewChecker(ctx, stock)
	defects := checker.CheckFundamentals(ctx)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	desc := "OK"
	descColor := tablewriter.FgHiGreenColor
	if len(defects) > 0 {
		table.SetHeader([]string{"FAILED ITEM", "FAILED REASON"})
		table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold})
		desc = "FAILED"
		descColor = tablewriter.FgHiRedColor
	}
	table.SetFooter([]string{fmt.Sprintf("%s %s", result.Name, result.Secucode), desc}) // Add Footer
	table.SetFooterColor(tablewriter.Colors{}, tablewriter.Colors{tablewriter.Bold, descColor})

	table.AppendBulk(defects) // Add Bulk Data
	table.Render()
}
