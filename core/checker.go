// 对给定股票进行检测

package core

import (
	"context"
	"fmt"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/x-stock/datacenter/eastmoney"
	"github.com/axiaoxin-com/x-stock/model"
)

// CheckerOptions 检测条件选项
type CheckerOptions struct {
	// ROE 不低于该值
	MinROE float64 `json:"min_roe"`
	// 连续增长年数
	CheckYears int `json:"check_years"`
	// ROE 高于该值时不做连续增长检查
	NoCheckYearsROE float64 `json:"no_check_years_roe"`
	// 最大资产负债率百分比(%)
	MaxDebtAssetRatio float64 `json:"max_debt_ratio"`
	// 最小历史波动率
	MinHV float64 `json:"min_hv"`
	// 最小市值（亿）
	MinTotalMarketCap float64 `json:"min_total_market_cap"`
}

// DefaultCheckerOptions 默认检测值
var DefaultCheckerOptions = CheckerOptions{
	MinROE:            8.0,
	CheckYears:        3,
	NoCheckYearsROE:   16.0,
	MaxDebtAssetRatio: 60.0,
	MinHV:             0.0,
	MinTotalMarketCap: 20.0,
}

// Checker 检测器实例
type Checker struct {
	Stock model.Stock
}

// NewChecker 创建检查器实例
func NewChecker(ctx context.Context, stock model.Stock) Checker {
	return Checker{
		Stock: stock,
	}
}

// CheckFundamentalsWithOptions 按条件检测股票基本面
func (c Checker) CheckFundamentalsWithOptions(ctx context.Context, options CheckerOptions) (defects [][]string) {
	// 最新 ROE 高于 n%
	if c.Stock.BaseInfo.RoeWeight < options.MinROE {
		checkItemName := "ROE VALUE"
		defect := fmt.Sprintf(
			"Latest ROE:%v is not greater than:%+v",
			c.Stock.BaseInfo.RoeWeight,
			options.MinROE,
		)
		defects = append(defects, []string{checkItemName, defect})
	}

	// ROE 均值小于 NoCheckYearsROE 时，至少 n 年内逐年递增
	roeavg, err := goutils.AvgFloat64(c.Stock.HistoricalFinaMainData.ROEList(ctx, options.CheckYears))
	if err != nil {
		logging.Warn(ctx, "roe avg error:"+err.Error())
	}
	if roeavg < options.NoCheckYearsROE &&
		!c.Stock.HistoricalFinaMainData.IsIncreasingByYears(ctx, "ROE", options.CheckYears) {
		checkItemName := "ROE INCREASING"
		defect := fmt.Sprintf(
			"ROE is not increasing in %d years. fina:%+v",
			options.CheckYears,
			c.Stock.HistoricalFinaMainData.ROEList(ctx, options.CheckYears),
		)
		defects = append(defects, []string{checkItemName, defect})
	}

	// EPS 至少 n 年内逐年递增
	if !c.Stock.HistoricalFinaMainData.IsIncreasingByYears(ctx, "EPS", options.CheckYears) {
		checkItemName := "EPS INCREASING"
		defect := fmt.Sprintf(
			"EPS is not increasing in %d years. fina:%+v",
			options.CheckYears,
			c.Stock.HistoricalFinaMainData.EPSList(ctx, options.CheckYears),
		)
		defects = append(defects, []string{checkItemName, defect})
	}

	// 营业总收入至少 n 年内逐年递增
	if !c.Stock.HistoricalFinaMainData.IsIncreasingByYears(ctx, "REVENUE", options.CheckYears) {
		checkItemName := "REVENUE INCREASING"
		defect := fmt.Sprintf(
			"REVENUE is not increasing in %d years. fina:%+v",
			options.CheckYears,
			c.Stock.HistoricalFinaMainData.RevenueList(ctx, options.CheckYears),
		)
		defects = append(defects, []string{checkItemName, defect})
	}

	// 净利润至少 n 年内逐年递增
	if !c.Stock.HistoricalFinaMainData.IsIncreasingByYears(ctx, "PROFIT", options.CheckYears) {
		checkItemName := "PROFIT INCREASING"
		defect := fmt.Sprintf(
			"PROFIT is not increasing in %d years. fina:%+v",
			options.CheckYears,
			c.Stock.HistoricalFinaMainData.ProfitList(ctx, options.CheckYears),
		)
		defects = append(defects, []string{checkItemName, defect})
	}

	// 估值较低或中等
	if c.Stock.ValuationStatus == eastmoney.ValuationHigh {
		checkItemName := "VALUATION STATUS"
		defect := "Valuation status is high"
		defects = append(defects, []string{checkItemName, defect})
	}

	// 股价低于合理价格
	if c.Stock.RightPrice != -1 && c.Stock.GetPrice() > c.Stock.RightPrice {
		checkItemName := "PRICE"
		defect := fmt.Sprintf(
			"NewPrice:%f is higher than RightPrice:%f",
			c.Stock.BaseInfo.NewPrice,
			c.Stock.RightPrice,
		)
		defects = append(defects, []string{checkItemName, defect})
	}

	// 负债率低于 MaxDebtRatio （可选条件）
	if options.MaxDebtAssetRatio != 0 {
		checkItemName := "DEBT ASSET RATIO"
		if len(c.Stock.HistoricalFinaMainData) > 0 && c.Stock.HistoricalFinaMainData[0].Zcfzl > options.MaxDebtAssetRatio {
			defect := fmt.Sprintf(
				"DebtAssetRatio(Zcfzl):%f is greater than %f",
				c.Stock.HistoricalFinaMainData[0].Zcfzl,
				options.MaxDebtAssetRatio,
			)
			defects = append(defects, []string{checkItemName, defect})
		}
	}

	// 历史波动率 （可选条件）
	if options.MinHV != 0 {
		checkItemName := "HV"
		if c.Stock.HistoricalVolatility > options.MinHV {
			defect := fmt.Sprintf(
				"HistoricalVolatility:%f is greater than %f",
				c.Stock.HistoricalVolatility,
				options.MinHV,
			)
			defects = append(defects, []string{checkItemName, defect})
		}
	}

	// 市值
	if c.Stock.BaseInfo.TotalMarketCap < options.MinTotalMarketCap {
		checkItemName := "TOTAL MARKET CAP"
		defect := fmt.Sprintf(
			"TotalMarketCap:%f is less than %f",
			c.Stock.BaseInfo.TotalMarketCap,
			options.MinTotalMarketCap,
		)
		defects = append(defects, []string{checkItemName, defect})
	}
	return
}

// CheckFundamentals 按默认条件进行基本面检测
func (c Checker) CheckFundamentals(ctx context.Context) [][]string {
	return c.CheckFundamentalsWithOptions(ctx, DefaultCheckerOptions)
}
