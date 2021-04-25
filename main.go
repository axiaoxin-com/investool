// Package main x-stock is my stock bot
package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/x-stock/checker"
	"github.com/axiaoxin-com/x-stock/exportor"
	"github.com/axiaoxin-com/x-stock/ui/terminal"
)

const (
	// ProcessorExportor 导出器
	ProcessorExportor = "exportor"
	// ProcessorChecker 检测器
	ProcessorChecker = "checker"
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
	// 要检测的关键词
	checkKeyword string
)

// 解析命令行启动参数
func parseFlags() {

	flag.StringVar(&processor, "run", "", "set processor to run: exportor|checker|tview|webserver")
	flag.StringVar(&loglevel, "l", "info", "set log level: debug|info|warn|error")
	flag.StringVar(
		&exportFilename,
		"f",
		fmt.Sprintf("./docs/x-stock.%s.xlsx", time.Now().Format("20060102")),
		"exportor arg to set filename, support .json .csv .xlsx .png",
	)
	flag.StringVar(&checkKeyword, "k", "", "checker arg set check keyword: <name>|<code>")
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
	case ProcessorChecker:
		ctx := context.Background()
		checker.Check(ctx, checkKeyword)
	case ProcessorTview:
		terminal.Run()
	case ProcessorWebserver:
	default:
		flag.PrintDefaults()
	}
}
