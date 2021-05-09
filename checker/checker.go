// Package checker 对给定股票名/股票代码进行检测
package checker

import (
	"context"
	"fmt"
	"os"

	"github.com/axiaoxin-com/x-stock/core"
	"github.com/olekukonko/tablewriter"
)

// Check 对给定名称或代码进行检测，输出检测结果
func Check(ctx context.Context, keywords []string, opts core.CheckerOptions) (map[string][][]string, error) {
	searcher := core.NewSearcher(ctx)
	stocks, err := searcher.Search(ctx, keywords)
	if err != nil {
		table := newTable()
		data := []string{"内部错误", err.Error()}
		renderTable(table, [][]string{data}, []string{"", "ERROR"})
		return nil, err
	}
	results := map[string][][]string{}
	for _, stock := range stocks {
		checker := core.NewChecker(ctx, opts)
		defects := checker.CheckFundamentals(ctx, stock)
		table := newTable()
		k := fmt.Sprintf("%s-%s", stock.BaseInfo.SecurityNameAbbr, stock.BaseInfo.Secucode)
		if len(defects) > 0 {
			results[k] = defects
			renderTable(table, defects, []string{k, "FAILED"})
		} else {
			data := []string{"--", "--"}
			renderTable(table, [][]string{data}, []string{k, "OK"})
		}
	}
	return results, nil
}

func newTable() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetRowLine(true)
	headers := []string{"未通过检测的指标", "原因"}
	table.SetHeader(headers)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold})
	return table
}

func renderTable(table *tablewriter.Table, data [][]string, footers []string) {
	footerValColor := tablewriter.FgHiRedColor
	if footers[1] == "OK" {
		footerValColor = tablewriter.FgHiGreenColor
	}
	table.SetFooter(footers)
	table.SetFooterColor(tablewriter.Colors{}, tablewriter.Colors{tablewriter.Bold, footerValColor})
	table.AppendBulk(data)
	table.Render()
}
