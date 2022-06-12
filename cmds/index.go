package cmds

import (
	"fmt"
	"os"
	"strconv"

	"github.com/axiaoxin-com/investool/datacenter/eastmoney"
	"github.com/axiaoxin-com/logging"
	"github.com/olekukonko/tablewriter"
)

func showIndexData(data *eastmoney.IndexData) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetRowSeparator("")
	table.SetBorder(false)
	table.SetNoWhiteSpace(true)
	headers := []string{}
	table.SetHeader(headers)
	table.SetCaption(true, "指数信息")
	rows := [][]string{
		{"指数名称", data.FullIndexName},
		{"指数说明", data.Reaprofile},
		{"编制方", data.MakerName},
		{"估值", data.IndexValueCN()},
		{"估值PE值", data.Petim},
		{"估值PE百分位", data.Pep100},
		{"指数代码", data.IndexCode},
		{"板块名称", data.BKName},
		{"当前点数", data.NewPrice},
		{"最新涨幅", data.NewCHG},
		{"最近一周涨幅", data.W},
		{"最近一月涨幅", data.M},
		{"最近三月涨幅", data.Q},
		{"最近六月涨幅", data.Hy},
		{"最近一年涨幅", data.Y},
		{"最近两年涨幅", data.Twy},
		{"最近三年涨幅", data.Try},
		{"最近五年涨幅", data.Fy},
		{"今年来涨幅", data.Sy},
	}
	table.AppendBulk(rows)

	table.Render()
}
func showIndexStocks(stocks []eastmoney.ZSCFGItem) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetRowLine(true)
	table.SetCaption(true, "指数成分股")
	headers := []string{"股票名称", "股票代码", "持仓占比"}
	table.SetHeader(headers)

	sum := 0.0
	for _, stock := range stocks {
		row := []string{stock.StockName, stock.StockCode, stock.Marketcappct}
		table.Append(row)
		if stock.Marketcappct != "" && stock.Marketcappct != "--" {
			v, err := strconv.ParseFloat(stock.Marketcappct, 64)
			if err != nil {
				logging.Error(nil, err.Error())
				continue
			}
			sum += v
		}
	}

	footers := []string{fmt.Sprintf("总数:%d", len(stocks)), "--", fmt.Sprintf("占比求和:%.2f", sum)}
	table.SetFooter(footers)
	table.Render()
}
