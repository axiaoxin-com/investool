// Package exportor 导出各类型的数据结果
package exportor

import (
	"context"
	"os"
	"path"
	"strings"
	"time"

	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/x-stock/model"
	"github.com/axiaoxin-com/x-stock/parser"
)

// Exportor exportor 实例
type Exportor struct {
	Stocks        DataList
	FilterOptions parser.FilterOptions
}

// New 创建要导出的数据列表
func New(ctx context.Context, stocks model.StockList, filterOptions parser.FilterOptions) Exportor {
	dlist := DataList{}
	for _, s := range stocks {
		dlist = append(dlist, NewData(s))
	}

	return Exportor{
		Stocks:        dlist,
		FilterOptions: filterOptions,
	}
}

// Export 导出数据
func Export(ctx context.Context, exportFilename string) {
	beginTime := time.Now()
	filedir := path.Dir(exportFilename)
	fileext := strings.ToLower(path.Ext(exportFilename))
	exportType := "csv"
	switch fileext {
	case ".json":
		exportType = "json"
	case ".csv", ".txt":
		exportType = "csv"
	case ".xlsx", ".xls":
		exportType = "excel"
	case ".png", ".jpg", ".jpeg", ".pic":
		exportType = "pic"
	}
	if _, err := os.Stat(filedir); os.IsNotExist(err) {
		os.Mkdir(filedir, 0755)
	}

	logging.Infof(ctx, "x-stock exportor start export selected stocks to %s", exportFilename)
	var err error
	stocks, err := parser.AutoFilterStocks(ctx)
	if err != nil {
		logging.Fatal(ctx, err.Error())
	}
	e := New(ctx, stocks, parser.DefaultFilterOptions)
	e.Stocks.SortByROE()

	switch exportType {
	case "json":
		_, err = e.ExportJSON(ctx, exportFilename)
	case "csv":
		_, err = e.ExportCSV(ctx, exportFilename)
	case "excel":
		_, err = e.ExportExcel(ctx, exportFilename)
	case "pic":
		// _, err = e.ExportPic(ctx, exportFilename)
	}
	if err != nil {
		logging.Fatal(ctx, err.Error())
	}

	latency := time.Now().Sub(beginTime).Milliseconds()
	logging.Infof(ctx, "x-stock exportor export %s succuss, latency:%+v", exportType, latency)
}
