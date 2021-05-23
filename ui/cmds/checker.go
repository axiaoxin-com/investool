// 检测器 cli command

package cmds

import (
	"context"
	"strings"

	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/x-stock/checker"
	"github.com/axiaoxin-com/x-stock/core"
	"github.com/urfave/cli/v2"
)

const (
	// ProcessorChecker 检测器
	ProcessorChecker = "checker"
)

// FlagsChecker cli flags
func FlagsChecker() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "keyword",
			Aliases:  []string{"k"},
			Value:    "",
			Usage:    "检给定股票名称或代码，多个股票批量检测使用/分割。如: 招商银行/中国平安/600519",
			Required: true,
		},
	}
}

// FlagsCheckerOptions exportor checker flags
func FlagsCheckerOptions() []cli.Flag {
	return []cli.Flag{
		&cli.Float64Flag{
			Name:        "checker.min_roe",
			Value:       8.0,
			Usage:       "最新一期 ROE 不低于该值",
			DefaultText: "8.0",
		},
		&cli.IntFlag{
			Name:        "checker.check_years",
			Value:       5,
			Usage:       "连续增长年数",
			DefaultText: "5",
		},
		&cli.Float64Flag{
			Name:        "checker.no_check_years_roe",
			Value:       20.0,
			Usage:       "ROE 高于该值时不做连续增长检查",
			DefaultText: "20.0",
		},
		&cli.Float64Flag{
			Name:        "checker.max_debt_asset_ratio",
			Value:       60.0,
			Usage:       "最大资产负债率百分比(%)",
			DefaultText: "60.0",
		},
		&cli.Float64Flag{
			Name:        "checker.min_hv",
			Value:       1.0,
			Usage:       "最大历史波动率",
			DefaultText: "1.0",
		},
		&cli.Float64Flag{
			Name:        "checker.min_total_market_cap",
			Value:       100.0,
			Usage:       "最小市值（亿）",
			DefaultText: "100.0",
		},
		&cli.Float64Flag{
			Name:        "checker.bank_min_roa",
			Value:       0.5,
			Usage:       "银行股最小 ROA",
			DefaultText: "0.5",
		},
		&cli.Float64Flag{
			Name:        "checker.bank_min_zbczl",
			Value:       8.0,
			Usage:       "银行股最小资本充足率",
			DefaultText: "8.0",
		},
		&cli.Float64Flag{
			Name:        "checker.bank_max_bldkl",
			Value:       3.0,
			Usage:       "银行股最大不良贷款率",
			DefaultText: "3.0",
		},
		&cli.Float64Flag{
			Name:        "checker.bank_min_bldkbbfgl",
			Value:       100.0,
			Usage:       "银行股最低不良贷款拨备覆盖率",
			DefaultText: "100.0",
		},
		&cli.BoolFlag{
			Name:        "checker.is_check_mll_stability",
			Value:       false,
			Usage:       "是否检测毛利率稳定性",
			DefaultText: "false",
		},
		&cli.BoolFlag{
			Name:        "checker.is_check_jll_stability",
			Value:       false,
			Usage:       "是否检测净利率稳定性",
			DefaultText: "false",
		},
		&cli.BoolFlag{
			Name:        "checker.is_check_price_by_calc",
			Value:       true,
			Usage:       "是否使用估算合理价进行检测，高于估算价将被过滤",
			DefaultText: "true",
		},
		&cli.Float64Flag{
			Name:        "checker.max_peg",
			Value:       1.5,
			Usage:       "最大 PEG",
			DefaultText: "1.5",
		},
		&cli.Float64Flag{
			Name:        "checker.min_byys_ratio",
			Value:       0.9,
			Usage:       "最小本业营收比",
			DefaultText: "0.9",
		},
		&cli.Float64Flag{
			Name:        "checker.max_byys_ratio",
			Value:       1.1,
			Usage:       "最大本业营收比",
			DefaultText: "1.1",
		},
		&cli.Float64Flag{
			Name:        "checker.min_fzldb",
			Value:       1.0,
			Usage:       "最小负债流动比",
			DefaultText: "1.0",
		},
	}
}

// NewCheckerOptions 从命令行参数解析 CheckerOptions
func NewCheckerOptions(c *cli.Context) core.CheckerOptions {
	checkerOpts := core.DefaultCheckerOptions
	checkerOpts.MinROE = c.Float64("checker.min_roe")
	checkerOpts.CheckYears = c.Int("checker.check_years")
	checkerOpts.NoCheckYearsROE = c.Float64("checker.no_check_years_roe")
	checkerOpts.MaxDebtAssetRatio = c.Float64("checker.max_debt_asset_ratio")
	checkerOpts.MaxHV = c.Float64("checker.max_hv")
	checkerOpts.MinTotalMarketCap = c.Float64("checker.min_total_market_cap")
	checkerOpts.BankMinROA = c.Float64("checker.bank_min_roa")
	checkerOpts.BankMinZBCZL = c.Float64("checker.bank_min_zbczl")
	checkerOpts.BankMaxBLDKL = c.Float64("checker.bank_max_bldkl")
	checkerOpts.BankMinBLDKBBFGL = c.Float64("checker.bank_min_bldkbbfgl")
	checkerOpts.IsCheckMLLStability = c.Bool("checker.is_check_mll_stability")
	checkerOpts.IsCheckJLLStability = c.Bool("checker.is_check_jll_stability")
	checkerOpts.IsCheckPriceByCalc = c.Bool("checker.is_check_price_by_calc")
	checkerOpts.MaxPEG = c.Float64("checker.max_peg")
	checkerOpts.MinBYYSRatio = c.Float64("checker.min_byys_ratio")
	checkerOpts.MaxBYYSRatio = c.Float64("checker.max_byys_ratio")
	checkerOpts.MinFZLDB = c.Float64("checker.min_fzldb")
	return checkerOpts
}

// ActionChecker cli action
func ActionChecker() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		loglevel := c.String("loglevel")
		logging.SetLevel(loglevel)
		keyword := c.String("keyword")
		ctx := context.Background()
		keywords := strings.Split(keyword, "/")
		opts := NewCheckerOptions(c)
		checker.Check(ctx, keywords, opts)
		return nil
	}
}

// CommandChecker 检测器 cli command
func CommandChecker() *cli.Command {
	flags := FlagsChecker()
	flags = append(flags, FlagsCheckerOptions()...)
	cmd := &cli.Command{
		Name:   ProcessorChecker,
		Usage:  "股票检测器",
		Flags:  flags,
		Action: ActionChecker(),
	}
	return cmd
}
