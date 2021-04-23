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
	Stocks DataList
}

// New 创建要导出的数据列表
func New(ctx context.Context, stocks model.StockList) Exportor {
	dlist := DataList{}
	for _, s := range stocks {
		dlist = append(dlist, NewData(s))
	}

	return Exportor{
		Stocks: dlist,
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
	case ".csv":
		exportType = "csv"
	}
	if _, err := os.Stat(filedir); os.IsNotExist(err) {
		os.Mkdir(filedir, 0755)
	}
	logging.Infof(ctx, "x-stock exportor start export selected stocks to %s", exportFilename)
	stocks, err := parser.AutoFilterStocks(ctx)
	if err != nil {
		logging.Fatal(ctx, err.Error())
	}
	e := New(ctx, stocks)
	e.Stocks.SortByROE()

	switch exportType {
	case "json":
		_, err := e.ExportJSON(ctx, exportFilename)
		if err != nil {
			logging.Fatal(ctx, err.Error())
		}
	case "csv":
		_, err := e.ExportCSV(ctx, exportFilename)
		if err != nil {
			logging.Fatal(ctx, err.Error())
		}
	}

	latency := time.Now().Sub(beginTime).Milliseconds()
	logging.Infof(ctx, "x-stock exportor export %s succuss, latency:%+v", exportType, latency)
}
