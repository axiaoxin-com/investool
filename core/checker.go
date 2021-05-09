// 对给定股票进行检测

package core

import (
	"context"
	"fmt"
	"strings"

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
	// 是否使用估算合理价进行检测，高于估算价将被过滤
	IsCheckPriceByCalc bool `json:"is_check_price_by_calc"`
	// 最大 PEG
	MaxPEG float64 `json:"max_peg"`
	// 最小本业营收比
	MinBYYSRatio float64 `json:"min_byys_ratio"`
	// 最大本业营收比
	MaxBYYSRatio float64 `json:"max_byys_ratio"`
}

// DefaultCheckerOptions 默认检测值
var DefaultCheckerOptions = CheckerOptions{
	MinROE:              8.0,
	CheckYears:          5,
	NoCheckYearsROE:     20.0,
	MaxDebtAssetRatio:   60.0,
	MaxHV:               1.0,
	MinTotalMarketCap:   100.0,
	BankMinROA:          0.5,
	BankMinZBCZL:        8.0,
	BankMaxBLDKL:        3.0,
	BankMinBLDKBBFGL:    100.0,
	IsCheckJLLStability: true,
	IsCheckMLLStability: true,
	IsCheckPriceByCalc:  false,
	MaxPEG:              1.5,
	MinBYYSRatio:        0.9,
	MaxBYYSRatio:        1.1,
}

// Checker 检测器实例
type Checker struct {
	Options CheckerOptions
}

// NewChecker 创建检查器实例
func NewChecker(ctx context.Context, opts CheckerOptions) Checker {
	return Checker{
		Options: opts,
	}
}

// CheckFundamentals 检测股票基本面
// [[检测失败项, 原因], ...]
func (c Checker) CheckFundamentals(ctx context.Context, stock model.Stock) (defects [][]string) {
	// ROE 高于 n%
	if stock.BaseInfo.RoeWeight < c.Options.MinROE {
		checkItemName := "净资产收益率 (ROE)"
		defect := fmt.Sprintf("最新一期 ROE:%f 低于:%f", stock.BaseInfo.RoeWeight, c.Options.MinROE)
		defects = append(defects, []string{checkItemName, defect})
	}

	// ROE 均值小于 NoCheckYearsROE 时，至少 n 年内逐年递增
	roeList := stock.HistoricalFinaMainData.ValueList(ctx, "ROE", c.Options.CheckYears)
	roeavg, err := goutils.AvgFloat64(roeList)
	if err != nil {
		logging.Warn(ctx, "roe avg error:"+err.Error())
	}
	if roeavg < c.Options.NoCheckYearsROE &&
		!stock.HistoricalFinaMainData.IsIncreasingByYears(ctx, "ROE", c.Options.CheckYears) {
		checkItemName := "ROE 逐年递增"
		defect := fmt.Sprintf("%d 年内未逐年递增:%+v", c.Options.CheckYears, roeList)
		defects = append(defects, []string{checkItemName, defect})
	}

	// EPS 至少 n 年内逐年递增
	if !stock.HistoricalFinaMainData.IsIncreasingByYears(ctx, "EPS", c.Options.CheckYears) {
		checkItemName := "EPS 逐年递增"
		defect := fmt.Sprintf(
			"%d 年内未逐年递增:%+v",
			c.Options.CheckYears,
			stock.HistoricalFinaMainData.ValueList(ctx, "EPS", c.Options.CheckYears),
		)
		defects = append(defects, []string{checkItemName, defect})
	}

	// 营业总收入至少 n 年内逐年递增
	if !stock.HistoricalFinaMainData.IsIncreasingByYears(ctx, "REVENUE", c.Options.CheckYears) {
		checkItemName := "营收逐年递增"
		defect := fmt.Sprintf(
			"%d 年内未逐年递增:%+v",
			c.Options.CheckYears,
			stock.HistoricalFinaMainData.ValueList(ctx, "REVENUE", c.Options.CheckYears),
		)
		defects = append(defects, []string{checkItemName, defect})
	}

	// 净利润至少 n 年内逐年递增
	if !stock.HistoricalFinaMainData.IsIncreasingByYears(ctx, "PROFIT", c.Options.CheckYears) {
		checkItemName := "净利润逐年递增"
		defect := fmt.Sprintf(
			"%d 年内未逐年递增:%+v",
			c.Options.CheckYears,
			stock.HistoricalFinaMainData.ValueList(ctx, "NETPROFIT", c.Options.CheckYears),
		)
		defects = append(defects, []string{checkItemName, defect})
	}

	// 整体质地
	if !goutils.IsStrInSlice(stock.JZPG.GetValueTotalScore(), []string{"优秀", "良好"}) {
		checkItemName := "整体质地"
		defect := stock.JZPG.GetValueTotalScore()
		defects = append(defects, []string{checkItemName, defect})
	}

	// 行业均值水平估值
	if stock.JZPG.GetValuationScore() == "高于行业均值水平" {
		checkItemName := "行业均值水平估值"
		defect := stock.JZPG.GetValuationScore()
		defects = append(defects, []string{checkItemName, defect})
	}

	// 市盈率、市净率、市现率、市销率全部估值较高
	highValuation := true
	highValuationDesc := []string{}
	for k, v := range stock.ValuationMap {
		if v != "估值较高" {
			highValuation = false
			break
		}
		highValuationDesc = append(highValuationDesc, k+v)
	}
	if highValuation {
		checkItemName := "四率估值全部较高"
		defect := strings.Join(highValuationDesc, "\n")
		defects = append(defects, []string{checkItemName, defect})
	}

	// 股价低于合理价格
	if c.Options.IsCheckPriceByCalc {
		if stock.RightPrice != -1 && stock.GetPrice() > stock.RightPrice {
			checkItemName := "股价"
			defect := fmt.Sprintf("最新股价:%f 高于合理价:%f", stock.BaseInfo.NewPrice, stock.RightPrice)
			defects = append(defects, []string{checkItemName, defect})
		}
	}

	// 负债率低于 MaxDebtRatio （可选条件），金融股不检测该项
	if !goutils.IsStrInSlice(stock.GetOrgType(), []string{"银行", "保险"}) {
		if c.Options.MaxDebtAssetRatio != 0 && len(stock.HistoricalFinaMainData) > 0 {
			if stock.HistoricalFinaMainData[0].Zcfzl > c.Options.MaxDebtAssetRatio {
				checkItemName := "负债率"
				defect := fmt.Sprintf("负债率:%f 高于:%f", stock.HistoricalFinaMainData[0].Zcfzl, c.Options.MaxDebtAssetRatio)
				defects = append(defects, []string{checkItemName, defect})
			}
		}
	}

	// 历史波动率 （可选条件）
	if c.Options.MaxHV != 0 && stock.HistoricalVolatility > c.Options.MaxHV {
		checkItemName := "历史波动率"
		defect := fmt.Sprintf("历史波动率:%f 高于:%f", stock.HistoricalVolatility, c.Options.MaxHV)
		defects = append(defects, []string{checkItemName, defect})

	}

	// 市值
	if stock.BaseInfo.TotalMarketCap < c.Options.MinTotalMarketCap {
		checkItemName := "市值"
		defect := fmt.Sprintf("市值:%f 低于:%f", stock.BaseInfo.TotalMarketCap, c.Options.MinTotalMarketCap)
		defects = append(defects, []string{checkItemName, defect})
	}

	// 银行股特殊检测
	if stock.GetOrgType() == "银行" && len(stock.HistoricalFinaMainData) > 0 {
		fmd := stock.HistoricalFinaMainData[0]
		if stock.BaseInfo.ROA < c.Options.BankMinROA {
			checkItemName := "总资产收益率 (ROA)"
			defect := fmt.Sprintf("ROA:%f 低于:%f", stock.BaseInfo.ROA, c.Options.BankMinROA)
			defects = append(defects, []string{checkItemName, defect})
		}
		if fmd.Newcapitalader < c.Options.BankMinZBCZL {
			checkItemName := "资本充足率"
			defect := fmt.Sprintf("资本充足率:%f 低于:%f", fmd.Newcapitalader, c.Options.BankMinZBCZL)
			defects = append(defects, []string{checkItemName, defect})
		}
		if c.Options.BankMaxBLDKL != 0 && fmd.NonPerLoan > c.Options.BankMaxBLDKL {
			checkItemName := "不良贷款率"
			defect := fmt.Sprintf("不良贷款率:%f 高于:%f", fmd.Newcapitalader, c.Options.BankMinZBCZL)
			defects = append(defects, []string{checkItemName, defect})
		}
		if fmd.Bldkbbl < c.Options.BankMinBLDKBBFGL {
			checkItemName := "不良贷款拨备覆盖率"
			defect := fmt.Sprintf("不良贷款拨备覆盖率:%f 低于:%f", fmd.Newcapitalader, c.Options.BankMinZBCZL)
			defects = append(defects, []string{checkItemName, defect})
		}
	}

	// 毛利率稳定性 （只检测非金融股）
	if c.Options.IsCheckMLLStability && !goutils.IsStrInSlice(stock.GetOrgType(), []string{"银行", "保险"}) {
		if stock.HistoricalFinaMainData.IsStability(ctx, "MLL", c.Options.CheckYears) {
			checkItemName := "毛利率稳定性"
			defect := fmt.Sprintf(
				"%d 年内稳定性较差:%v",
				c.Options.CheckYears,
				stock.HistoricalFinaMainData.ValueList(ctx, "MLL", c.Options.CheckYears),
			)
			defects = append(defects, []string{checkItemName, defect})
		}
	}

	// 净利率稳定性
	if c.Options.IsCheckMLLStability && stock.HistoricalFinaMainData.IsStability(ctx, "JLL", c.Options.CheckYears) {
		checkItemName := "净利率稳定性"
		defect := fmt.Sprintf(
			"%d 年内稳定性较差:%v",
			c.Options.CheckYears,
			stock.HistoricalFinaMainData.ValueList(ctx, "JLL", c.Options.CheckYears),
		)
		defects = append(defects, []string{checkItemName, defect})
	}

	// PEG
	if c.Options.MaxPEG != 0 && stock.PEG > c.Options.MaxPEG {
		checkItemName := "PEG"
		defect := fmt.Sprintf("PEG:%v 高于:%v", stock.PEG, c.Options.MaxPEG)
		defects = append(defects, []string{checkItemName, defect})
	}

	// 本业营收比
	if c.Options.MinBYYSRatio != 0 && c.Options.MaxBYYSRatio != 0 {
		if stock.BYYSRatio > c.Options.MaxBYYSRatio || stock.BYYSRatio < c.Options.MinBYYSRatio {
			checkItemName := "本业营收比"
			defect := fmt.Sprintf("当前本业营收比:%v 超出范围:%v-%v", stock.BYYSRatio, c.Options.MinBYYSRatio, c.Options.MaxBYYSRatio)
			defects = append(defects, []string{checkItemName, defect})
		}
	}

	// 审计意见
	if stock.FinaReportOpinion != nil {
		opinion := stock.FinaReportOpinion.(string)
		if opinion != "标准无保留意见" {
			checkItemName := "财报审计意见"
			defect := opinion
			defects = append(defects, []string{checkItemName, defect})
		}
	}

	return
}
