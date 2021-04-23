// Package main x-stock is my stock bot
package main

import (
	"context"
	"flag"
	"strings"

	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/x-stock/exportor"
	"github.com/axiaoxin-com/x-stock/parser"
)

const (
	// ProcessorExportor 导出数据
	ProcessorExportor = "exportor"
	// ProcessorTview 运行终端界面
	ProcessorTview = "tview"
	// ProcessorWebserver 运行 web 服务器
	ProcessorWebserver = "webserver"
)

var (
	// 要启动运行的进程
	processor string
	// 日志级别
	loglevel string
	// 要导出的文件名
	exportFilename string
	// 要导出的数据类型
	exportType string
)

// 解析命令行启动参数
func parseFlags() {
	flag.StringVar(&processor, "run", "", "processor: exportor|tview|webserver")
	flag.StringVar(&loglevel, "l", "info", "loglevel: debug|info|warn|error")
	flag.StringVar(&exportFilename, "f", "./x-stock.csv", "export filename")
	flag.StringVar(&exportType, "t", "csv", "export data type: csv|json")

	flag.Parse()
}

func init() {
	logging.SetLevel(loglevel)
	parseFlags()
}

// RunExportor 运行 exportor 导出数据
func RunExportor(ctx context.Context) {
	stocks, err := parser.AutoFilterStocks(ctx)
	data := exportor.New(ctx, stocks)
	data.SortByROE()
	if err != nil {
		logging.Fatal(ctx, err.Error())
	}
	switch strings.ToLower(exportType) {
	case "json":
		c, err := data.ExportJSON(ctx, exportFilename)
		if err != nil {
			logging.Fatal(ctx, err.Error())
		}
		logging.Info(ctx, "json content:"+string(c))
	case "csv":
		c, err := data.ExportCSV(ctx, exportFilename)
		if err != nil {
			logging.Fatal(ctx, err.Error())
		}
		logging.Info(ctx, "csv content:"+string(c))
	}
}

func main() {
	switch processor {
	case ProcessorExportor:
		ctx := context.Background()
		RunExportor(ctx)
	case ProcessorTview:
	case ProcessorWebserver:
	default:
		flag.PrintDefaults()
	}
}
