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
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetRowLine(true)
	desc := "OK"
	descColor := tablewriter.FgHiGreenColor

	results, err := datacenter.QQ.KeywordSearch(ctx, keyword)
	if err != nil {
		renderFail(table, err.Error())
		return
	}
	if len(results) == 0 {
		renderFail(table, "无对应股票")
		return
	}
	result := results[0]
	logging.Infof(ctx, "search result:%+v", result)
	filter := eastmoney.DefaultFilter
	filter.SpecialSecurityCode = result.SecurityCode
	stocks, err := datacenter.EastMoney.QuerySelectedStocksWithFilter(ctx, filter)
	if err != nil {
		renderFail(table, err.Error())
		return
	}
	if len(stocks) == 0 {
		renderFail(table, "无股票数据")
		return
	}
	stock, err := model.NewStock(ctx, stocks[0], false)
	if err != nil {
		renderFail(table, err.Error())
		return
	}
	checker := core.NewChecker(ctx, stock)
	defects := checker.CheckFundamentals(ctx)
	if len(defects) > 0 {
		table.SetHeader([]string{"未通过检测的指标", "原因"})
		table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold})
		desc = "FAILED"
		descColor = tablewriter.FgHiRedColor
	}
	table.SetFooter([]string{fmt.Sprintf("%s %s", result.Name, result.Secucode), desc}) // Add Footer
	table.SetFooterColor(tablewriter.Colors{}, tablewriter.Colors{tablewriter.Bold, descColor})

	table.AppendBulk(defects) // Add Bulk Data
	table.Render()
}

func renderFail(table *tablewriter.Table, desc string) {
	descColor := tablewriter.FgHiRedColor
	table.SetFooter([]string{"", desc})
	table.SetFooterColor(tablewriter.Colors{}, tablewriter.Colors{tablewriter.Bold, descColor})
	table.Render()
}
