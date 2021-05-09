// Package exportor 导出各类型的数据结果
package exportor

import (
	"context"
	"os"
	"path"
	"strings"
	"time"

	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/x-stock/core"
	"github.com/axiaoxin-com/x-stock/datacenter/eastmoney"
	"github.com/axiaoxin-com/x-stock/model"
)

// Descriptions 数据备注信息
type Descriptions struct {
	Filter         eastmoney.Filter    `json:"filter"`
	CheckerOptions core.CheckerOptions `json:"checker_options"`
}

// Exportor exportor 实例
type Exportor struct {
	Stocks       DataList
	Descriptions Descriptions
}

// New 创建要导出的数据列表
func New(ctx context.Context, stocks model.StockList, filter eastmoney.Filter, checkerOptions core.CheckerOptions) Exportor {
	dlist := DataList{}
	for _, s := range stocks {
		dlist = append(dlist, NewData(ctx, s))
	}

	return Exportor{
		Stocks: dlist,
		Descriptions: Descriptions{
			Filter:         filter,
			CheckerOptions: checkerOptions,
		},
	}
}

// Export 导出数据
func Export(ctx context.Context, exportFilename string, disableCheck bool) {
	beginTime := time.Now()
	filedir := path.Dir(exportFilename)
	fileext := strings.ToLower(path.Ext(exportFilename))
	exportType := "excel"
	switch fileext {
	case ".json":
		exportType = "json"
	case ".csv", ".txt":
		exportType = "csv"
	case ".xlsx", ".xls":
		exportType = "excel"
	case ".png", ".jpg", ".jpeg", ".pic":
		exportType = "pic"
	case ".all":
		exportType = "all"
	}
	if _, err := os.Stat(filedir); os.IsNotExist(err) {
		os.Mkdir(filedir, 0755)
	}

	logging.Infof(ctx, "x-stock exportor start export selected stocks to %s", exportFilename)
	var err error
	// 自动筛选股票
	selector := core.NewSelector(ctx)
	filter := eastmoney.DefaultFilter
	stocks, err := selector.AutoFilterStocksWithFilter(ctx, filter, disableCheck)
	if err != nil {
		logging.Fatal(ctx, err.Error())
	}
	checkOpt := core.DefaultCheckerOptions
	e := New(ctx, stocks, filter, checkOpt)

	switch exportType {
	case "json":
		_, err = e.ExportJSON(ctx, exportFilename)
	case "csv":
		_, err = e.ExportCSV(ctx, exportFilename)
	case "excel":
		_, err = e.ExportExcel(ctx, exportFilename)
	case "pic":
		_, err = e.ExportPic(ctx, exportFilename)
	case "all":
		jsonFilename := strings.ReplaceAll(exportFilename, ".all", ".json")
		_, err = e.ExportJSON(ctx, jsonFilename)
		csvFilename := strings.ReplaceAll(exportFilename, ".all", ".csv")
		_, err = e.ExportCSV(ctx, csvFilename)
		xlsxFilename := strings.ReplaceAll(exportFilename, ".all", ".xlsx")
		_, err = e.ExportExcel(ctx, xlsxFilename)
		pngFilename := strings.ReplaceAll(exportFilename, ".all", ".png")
		_, err = e.ExportPic(ctx, pngFilename)
	}
	if err != nil {
		logging.Fatal(ctx, err.Error())
	}

	logging.Infof(ctx, "x-stock exportor export %s succuss, latency:%#v", exportType, time.Now().Sub(beginTime))
}
