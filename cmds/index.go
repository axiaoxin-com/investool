package cmds

import (
	"fmt"
	"os"
	"strconv"

	"github.com/axiaoxin-com/investool/datacenter/eastmoney"
	"github.com/axiaoxin-com/logging"
	"github.com/olekukonko/tablewriter"
)

func showIndexStocks(stocks []eastmoney.ZSCFGItem) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetRowLine(true)
	headers := []string{"股票名称", "股票代码", "持仓占比"}
	table.SetHeader(headers)
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.BgBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.BgBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.BgBlackColor},
	)

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
