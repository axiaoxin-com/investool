// 指数数据

package cmds

import (
	"context"

	"github.com/axiaoxin-com/investool/datacenter"
	"github.com/urfave/cli/v2"
)

const (
	// ProcessorIndex 指数数据处理
	ProcessorIndex = "index"
)

// FlagsIndex cli flags
func FlagsIndex() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "index",
			Aliases:  []string{"i"},
			Value:    "",
			Usage:    "指定指数代码",
			Required: true,
		},
		&cli.BoolFlag{
			Name:     "stocks",
			Aliases:  []string{"s"},
			Value:    false,
			Usage:    "返回指数成分股",
			Required: false,
		},
	}
}

// ActionIndex cli action
func ActionIndex() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		ctx := context.Background()
		indexCode := c.String("index")
		showStocks := c.Bool("stocks")
		if showStocks {
			stocks, err := datacenter.EastMoney.ZSCFG(ctx, indexCode)
			if err != nil {
				return err
			}
			showIndexStocks(stocks)
		}
		return nil
	}
}

// CommandIndex 指数成分股 cli command
func CommandIndex() *cli.Command {
	flags := FlagsIndex()
	cmd := &cli.Command{
		Name:   ProcessorIndex,
		Usage:  "指数数据",
		Flags:  flags,
		Action: ActionIndex(),
	}
	return cmd
}
