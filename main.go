// Package main x-stock is my stock bot
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/x-stock/ui/cmds"
	"github.com/urfave/cli/v2"
)

// VERSION 版本号
const VERSION = "1.0.5"

var (
	// DefaultLoglevel 日志级别默认值
	DefaultLoglevel = "info"
	// ProcessorOptions 要启动运行的进程可选项
	ProcessorOptions = []string{cmds.ProcessorChecker, cmds.ProcessorExportor}
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
	app.BashComplete = func(c *cli.Context) {
		if c.NArg() > 0 {
			return
		}
		for _, i := range ProcessorOptions {
			fmt.Println(i)
		}
	}

	app.Commands = append(app.Commands, cmds.CommandExportor())
	app.Commands = append(app.Commands, cmds.CommandChecker())

	if err := app.Run(os.Args); err != nil {
		logging.Fatal(nil, err.Error())
	}
}
