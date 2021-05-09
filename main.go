// Package main x-stock is my stock bot
package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/x-stock/checker"
	"github.com/axiaoxin-com/x-stock/exportor"
	"github.com/urfave/cli/v2"
)

// VERSION 版本号
const VERSION = "0.0.6"

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
	// DefaultProcessor 要启动运行的进程默认值
	DefaultProcessor = ProcessorExportor
	// DefaultLoglevel 日志级别默认值
	DefaultLoglevel = "info"
	// DefaultExportFilename 要导出的文件名默认值
	DefaultExportFilename = fmt.Sprintf("./dist/x-stock.%s.xlsx", time.Now().Format("20060102"))

	// LogLevelOptions 日志级别参数的可选项
	LogLevelOptions = []string{"debug", "info", "warn", "error"}
	// ProcessorOptions 要启动运行的进程可选项
	ProcessorOptions = []string{ProcessorTview, ProcessorChecker, ProcessorExportor, ProcessorWebserver}
)

func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = "x-stock"
	app.Usage = "axiaoxin 的股票工具程序"
	app.UsageText = `该程序不构成任何投资建议，程序只是个人辅助工具，具体分析仍然需要自己判断。

官网地址: http://x-stock.axiaoxin.com`
	app.Version = VERSION
	app.Compiled = time.Now()
	app.Authors = []*cli.Author{
		{
			Name:  "axiaoxin",
			Email: "254606826@qq.com",
		},
	}
	app.Copyright = "(c) 2021 axiaoxin"

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "show the version",
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "loglevel",
			Aliases:     []string{"l"},
			Value:       DefaultLoglevel,
			Usage:       "日志级别 [debug|info|warn|error]",
			EnvVars:     []string{"XSTOCK_LOGLEVEL"},
			DefaultText: DefaultLoglevel,
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:      ProcessorExportor,
			Usage:     "股票筛选导出器",
			UsageText: "将按条件筛选出的股票导出到文件",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "filename",
					Aliases:     []string{"f"},
					Value:       DefaultExportFilename,
					Usage:       `导出筛选出的股票数据到指定文件，根据文件后缀名自动判断导出类型。支持的后缀名：[xlsx|csv|json|png|all]，all 表示导出全部支持的类型。`,
					EnvVars:     []string{"XSTOCK_EXPORTOR_FILENAME"},
					DefaultText: DefaultExportFilename,
				},
				&cli.BoolFlag{
					Name:        "disable_check",
					Aliases:     []string{"C"},
					Value:       false,
					Usage:       "关闭基本面检测，导出所有原始筛选结果",
					EnvVars:     []string{"XSTOCK_EXPORTOR_DISABLE_CHECK"},
					DefaultText: "false",
				},
			},
			Action: func(c *cli.Context) error {
				loglevel := c.String("loglevel")
				logging.SetLevel(loglevel)
				disableCheck := c.Bool("disable_check")
				ctx := context.Background()
				exportor.Export(ctx, c.String("filename"), disableCheck)
				return nil
			},
			BashComplete: func(c *cli.Context) {
				if c.NArg() > 0 {
					return
				}
				for _, i := range ProcessorOptions {
					fmt.Println(i)
				}
			},
		},
		{
			Name:  "checker",
			Usage: "股票检测器",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "keyword",
					Aliases:  []string{"k"},
					Value:    "",
					Usage:    "检给定股票名称或代码，多个股票批量检测使用/分割。如: 招商银行/中国平安/600519",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				loglevel := c.String("loglevel")
				logging.SetLevel(loglevel)
				keyword := c.String("keyword")
				ctx := context.Background()
				keywords := strings.Split(keyword, "/")
				checker.Check(ctx, keywords)
				return nil
			},
			BashComplete: func(c *cli.Context) {
				if c.NArg() > 0 {
					return
				}
				for _, i := range LogLevelOptions {
					fmt.Println(i)
				}
			},
		},
		{
			Name:  "tview",
			Usage: "终端界面",
			Action: func(c *cli.Context) error {
				loglevel := c.String("loglevel")
				logging.SetLevel(loglevel)
				return nil
			},
		},
		{
			Name:  "webserver",
			Usage: "WEB 服务器",
			Action: func(c *cli.Context) error {
				loglevel := c.String("loglevel")
				logging.SetLevel(loglevel)
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logging.Fatal(nil, err.Error())
	}
}
