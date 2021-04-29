// 对给定股票进行检测

package core

import (
	"context"
	"fmt"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/x-stock/model"
)

// CheckerOptions 检测条件选项
type CheckerOptions struct {
	// 最新一期 ROE 不低于该值
	MinROE float64 `json:"min_roe"`
	// 连续增长年数
	CheckYears int `json:"check_years"`
	// ROE 高于该值时不做连续增长检查
	NoCheckYearsROE float64 `json:"no_check_years_roe"`
	// 最大资产负债率百分比(%)
	MaxDebtAssetRatio float64 `json:"max_debt_ratio"`
	// 最大历史波动率
	MaxHV float64 `json:"min_hv"`
	// 最小市值（亿）
	MinTotalMarketCap float64 `json:"min_total_market_cap"`
	// 银行股最小 ROA
	BankMinROA float64 `json:"bank_min_roa"`
	// 银行股最小资本充足率
	BankMinZBCZL float64 `json:"bank_min_zbczl"`
	// 银行股最大不良贷款率
	BankMaxBLDKL float64 `json:"bank_max_bldkl"`
	// 银行股最低不良贷款拨备覆盖率
	BankMinBLDKBBFGL float64 `json:"bank_min_bldkbbfgl"`
	// 是否检测毛利率稳定性
	IsCheckMLLStability bool `json:"is_check_mll_stability"`
	// 是否检测净利率稳定性
	IsCheckJLLStability bool `json:"is_check_jll_stability"`
}

// DefaultCheckerOptions 默认检测值
var DefaultCheckerOptions = CheckerOptions{
	MinROE:              8.0,
	CheckYears:          3,
	NoCheckYearsROE:     16.0,
	MaxDebtAssetRatio:   60.0,
	MaxHV:               3.0,
	MinTotalMarketCap:   20.0,
	BankMinROA:          0.5,
	BankMinZBCZL:        8.0,
	BankMaxBLDKL:        3.0,
	BankMinBLDKBBFGL:    100.0,
	IsCheckJLLStability: true,
	IsCheckMLLStability: true,
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
// [[检测失败项, 原因], ...]
func (c Checker) CheckFundamentalsWithOptions(ctx context.Context, options CheckerOptions) (defects [][]string) {
	// ROE 高于 n%
	if c.Stock.BaseInfo.RoeWeight < options.MinROE {
		checkItemName := "净资产收益率 (ROE)"
		defect := fmt.Sprintf("最新一期 ROE:%f 低于:%f", c.Stock.BaseInfo.RoeWeight, options.MinROE)
		defects = append(defects, []string{checkItemName, defect})
	}

	// ROE 均值小于 NoCheckYearsROE 时，至少 n 年内逐年递增
	roeList := c.Stock.HistoricalFinaMainData.ValueList(ctx, "ROE", options.CheckYears)
	roeavg, err := goutils.AvgFloat64(roeList)
	if err != nil {
		logging.Warn(ctx, "roe avg error:"+err.Error())
	}
	if roeavg < options.NoCheckYearsROE &&
		!c.Stock.HistoricalFinaMainData.IsIncreasingByYears(ctx, "ROE", options.CheckYears) {
		checkItemName := "ROE 逐年递增"
		defect := fmt.Sprintf("%d 年内未逐年递增:%+v", options.CheckYears, roeList)
		defects = append(defects, []string{checkItemName, defect})
	}

	// EPS 至少 n 年内逐年递增
	if !c.Stock.HistoricalFinaMainData.IsIncreasingByYears(ctx, "EPS", options.CheckYears) {
		checkItemName := "EPS 逐年递增"
		defect := fmt.Sprintf(
			"%d 年内未逐年递增:%+v",
			options.CheckYears,
			c.Stock.HistoricalFinaMainData.ValueList(ctx, "EPS", options.CheckYears),
		)
		defects = append(defects, []string{checkItemName, defect})
	}

	// 营业总收入至少 n 年内逐年递增
	if !c.Stock.HistoricalFinaMainData.IsIncreasingByYears(ctx, "REVENUE", options.CheckYears) {
		checkItemName := "营收逐年递增"
		defect := fmt.Sprintf(
			"%d 年内未逐年递增:%+v",
			options.CheckYears,
			c.Stock.HistoricalFinaMainData.ValueList(ctx, "REVENUE", options.CheckYears),
		)
		defects = append(defects, []string{checkItemName, defect})
	}

	// 净利润至少 n 年内逐年递增
	if !c.Stock.HistoricalFinaMainData.IsIncreasingByYears(ctx, "PROFIT", options.CheckYears) {
		checkItemName := "净利润逐年递增"
		defect := fmt.Sprintf(
			"%d 年内未逐年递增:%+v",
			options.CheckYears,
			c.Stock.HistoricalFinaMainData.ValueList(ctx, "NETPROFIT", options.CheckYears),
		)
		defects = append(defects, []string{checkItemName, defect})
	}

	// TODO:公司总体质地良好或优秀

	// 股价低于合理价格
	if c.Stock.RightPrice != -1 && c.Stock.GetPrice() > c.Stock.RightPrice {
		checkItemName := "股价"
		defect := fmt.Sprintf("最新股价:%f 高于合理价:%f", c.Stock.BaseInfo.NewPrice, c.Stock.RightPrice)
		defects = append(defects, []string{checkItemName, defect})
	}

	// 负债率低于 MaxDebtRatio （可选条件），金融股不检测该项
	if !goutils.IsStrInSlice(c.Stock.GetOrgType(), []string{"银行", "保险"}) {
		if options.MaxDebtAssetRatio != 0 && len(c.Stock.HistoricalFinaMainData) > 0 {
			if c.Stock.HistoricalFinaMainData[0].Zcfzl > options.MaxDebtAssetRatio {
				checkItemName := "负债率"
				defect := fmt.Sprintf("负债率:%f 高于:%f", c.Stock.HistoricalFinaMainData[0].Zcfzl, options.MaxDebtAssetRatio)
				defects = append(defects, []string{checkItemName, defect})
			}
		}
	}

	// 历史波动率 （可选条件）
	if options.MaxHV != 0 && c.Stock.HistoricalVolatility > options.MaxHV {
		checkItemName := "历史波动率"
		defect := fmt.Sprintf("历史波动率:%f 高于:%f", c.Stock.HistoricalVolatility, options.MaxHV)
		defects = append(defects, []string{checkItemName, defect})

	}

	// 市值
	if c.Stock.BaseInfo.TotalMarketCap < options.MinTotalMarketCap {
		checkItemName := "市值"
		defect := fmt.Sprintf("市值:%f 低于:%f", c.Stock.BaseInfo.TotalMarketCap, options.MinTotalMarketCap)
		defects = append(defects, []string{checkItemName, defect})
	}

	// 银行股特殊检测
	if c.Stock.GetOrgType() == "银行" && len(c.Stock.HistoricalFinaMainData) > 0 {
		fmd := c.Stock.HistoricalFinaMainData[0]
		if c.Stock.BaseInfo.ROA < options.BankMinROA {
			checkItemName := "总资产收益率 (ROA)"
			defect := fmt.Sprintf("ROA:%f 低于:%f", c.Stock.BaseInfo.ROA, options.BankMinROA)
			defects = append(defects, []string{checkItemName, defect})
		}
		if fmd.Newcapitalader < options.BankMinZBCZL {
			checkItemName := "资本充足率"
			defect := fmt.Sprintf("资本充足率:%f 低于:%f", fmd.Newcapitalader, options.BankMinZBCZL)
			defects = append(defects, []string{checkItemName, defect})
		}
		if options.BankMaxBLDKL != 0 && fmd.NonPerLoan > options.BankMaxBLDKL {
			checkItemName := "不良贷款率"
			defect := fmt.Sprintf("不良贷款率:%f 高于:%f", fmd.Newcapitalader, options.BankMinZBCZL)
			defects = append(defects, []string{checkItemName, defect})
		}
		if fmd.Bldkbbl < options.BankMinBLDKBBFGL {
			checkItemName := "不良贷款拨备覆盖率"
			defect := fmt.Sprintf("不良贷款拨备覆盖率:%f 低于:%f", fmd.Newcapitalader, options.BankMinZBCZL)
			defects = append(defects, []string{checkItemName, defect})
		}
	}

	// 毛利率稳定性 （只检测非金融股）
	if options.IsCheckMLLStability && !goutils.IsStrInSlice(c.Stock.GetOrgType(), []string{"银行", "保险"}) {
		if c.Stock.HistoricalFinaMainData.IsStability(ctx, "MLL", options.CheckYears) {
			checkItemName := "毛利率稳定性"
			defect := fmt.Sprintf(
				"%d 年内稳定性较差:%v",
				options.CheckYears,
				c.Stock.HistoricalFinaMainData.ValueList(ctx, "MLL", options.CheckYears),
			)
			defects = append(defects, []string{checkItemName, defect})
		}
	}

	// 净利率稳定性
	if options.IsCheckMLLStability && c.Stock.HistoricalFinaMainData.IsStability(ctx, "JLL", options.CheckYears) {
		checkItemName := "净利率稳定性"
		defect := fmt.Sprintf(
			"%d 年内稳定性较差:%v",
			options.CheckYears,
			c.Stock.HistoricalFinaMainData.ValueList(ctx, "JLL", options.CheckYears),
		)
		defects = append(defects, []string{checkItemName, defect})
	}

	return
}

// CheckFundamentals 按默认条件进行基本面检测
func (c Checker) CheckFundamentals(ctx context.Context) [][]string {
	return c.CheckFundamentalsWithOptions(ctx, DefaultCheckerOptions)
}
