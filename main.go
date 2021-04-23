// Package main x-stock is my stock bot
package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/x-stock/exportor"
	"github.com/axiaoxin-com/x-stock/ui/terminal"
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
	flag.StringVar(
		&exportFilename,
		"f",
		fmt.Sprintf("./docs/x-stock.%s.xlsx", time.Now().Format("20060102")),
		"export filename",
	)
	flag.Parse()
}

func init() {
	parseFlags()
	logging.SetLevel(loglevel)
}

func main() {
	switch processor {
	case ProcessorExportor:
		ctx := context.Background()
		exportor.Export(ctx, exportFilename)
	case ProcessorTview:
		terminal.Run()
	case ProcessorWebserver:
	default:
		flag.PrintDefaults()
	}
}
