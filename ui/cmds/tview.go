// 终端界面 cli command

package cmds

import (
	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/x-stock/ui/terminal"
	"github.com/urfave/cli/v2"
)

const (
	// ProcessorTview 终端界面
	ProcessorTview = "tview"
)

// FlagsTview cli flags
func FlagsTview() []cli.Flag {
	return []cli.Flag{}
}

// ActionTview cli action
func ActionTview() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		loglevel := c.String("loglevel")
		logging.SetLevel(loglevel)
		return terminal.Run()
	}
}

// CommandTview 终端界面 cli command
func CommandTview() *cli.Command {
	cmd := &cli.Command{
		Name:   ProcessorTview,
		Usage:  "终端界面",
		Flags:  FlagsTview(),
		Action: ActionTview(),
	}
	return cmd
}
