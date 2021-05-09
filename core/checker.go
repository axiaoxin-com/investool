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
	MaxDebtAssetRatio float64 `json:"max_debt_asset_ratio"`
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

// CheckResult 检测结果
// key 为检测项，value为描述map {"ROE": {"desc": "高于8.0", "ok":"true"}}
type CheckResult map[string]map[string]string

// CheckFundamentals 检测股票基本面
// [[检测失败项, 原因], ...]
func (c Checker) CheckFundamentals(ctx context.Context, stock model.Stock) (result CheckResult, ok bool) {
	ok = true
	result = make(CheckResult)
	// ROE 高于 n%
	checkItemName := "净资产收益率 (ROE)"
	itemOK := true
	desc := fmt.Sprintf("最新一期 ROE:%f", stock.BaseInfo.RoeWeight)
	if stock.BaseInfo.RoeWeight < c.Options.MinROE {
		desc = fmt.Sprintf("最新一期 ROE:%f 低于:%f", stock.BaseInfo.RoeWeight, c.Options.MinROE)
		ok = false
		itemOK = false
	}
	result[checkItemName] = map[string]string{
		"desc": desc,
		"ok":   fmt.Sprint(itemOK),
	}

	// ROE 均值小于 NoCheckYearsROE 时，至少 n 年内逐年递增
	checkItemName = fmt.Sprintf("ROE 逐年递增（均值>=%f除外）", c.Options.NoCheckYearsROE)
	itemOK = true
	roeList := stock.HistoricalFinaMainData.ValueList(ctx, "ROE", c.Options.CheckYears)
	desc = fmt.Sprintf("%d 年内 ROE:%+v", c.Options.CheckYears, roeList)
	roeavg, err := goutils.AvgFloat64(roeList)
	if err != nil {
		logging.Warn(ctx, "roe avg error:"+err.Error())
	}
	if roeavg < c.Options.NoCheckYearsROE &&
		!stock.HistoricalFinaMainData.IsIncreasingByYears(ctx, "ROE", c.Options.CheckYears) {
		desc = fmt.Sprintf("%d 年内未逐年递增:%+v", c.Options.CheckYears, roeList)
		ok = false
		itemOK = false
	}
	result[checkItemName] = map[string]string{
		"desc": desc,
		"ok":   fmt.Sprint(itemOK),
	}

	// EPS 至少 n 年内逐年递增
	epsList := stock.HistoricalFinaMainData.ValueList(ctx, "EPS", c.Options.CheckYears)
	checkItemName = "EPS 逐年递增"
	itemOK = true
	desc = fmt.Sprintf("%d 年内 EPS:%+v", c.Options.CheckYears, epsList)
	if !stock.HistoricalFinaMainData.IsIncreasingByYears(ctx, "EPS", c.Options.CheckYears) {
		desc = fmt.Sprintf("%d 年内未逐年递增:%+v", c.Options.CheckYears, epsList)
		ok = false
		itemOK = false
	}
	result[checkItemName] = map[string]string{
		"desc": desc,
		"ok":   fmt.Sprint(itemOK),
	}

	// 营业总收入至少 n 年内逐年递增
	revList := stock.HistoricalFinaMainData.ValueList(ctx, "REVENUE", c.Options.CheckYears)
	checkItemName = "营收逐年递增"
	itemOK = true
	desc = fmt.Sprintf("%d 年内营收:%+v", c.Options.CheckYears, revList)
	if !stock.HistoricalFinaMainData.IsIncreasingByYears(ctx, "REVENUE", c.Options.CheckYears) {
		desc = fmt.Sprintf("%d 年内未逐年递增:%+v", c.Options.CheckYears, revList)
		ok = false
		itemOK = false
	}
	result[checkItemName] = map[string]string{
		"desc": desc,
		"ok":   fmt.Sprint(itemOK),
	}

	// 净利润至少 n 年内逐年递增
	netprofitList := stock.HistoricalFinaMainData.ValueList(ctx, "NETPROFIT", c.Options.CheckYears)
	checkItemName = "净利润逐年递增"
	itemOK = true
	desc = fmt.Sprintf("%d 年内净利润:%+v", c.Options.CheckYears, netprofitList)
	if !stock.HistoricalFinaMainData.IsIncreasingByYears(ctx, "PROFIT", c.Options.CheckYears) {
		desc = fmt.Sprintf("%d 年内未逐年递增:%+v", c.Options.CheckYears, netprofitList)
		ok = false
		itemOK = false
	}
	result[checkItemName] = map[string]string{
		"desc": desc,
		"ok":   fmt.Sprint(itemOK),
	}

	// 整体质地
	checkItemName = "整体质地"
	itemOK = true
	desc = stock.JZPG.GetValueTotalScore()
	if !goutils.IsStrInSlice(stock.JZPG.GetValueTotalScore(), []string{"优秀", "良好"}) {
		ok = false
		itemOK = false
	}
	result[checkItemName] = map[string]string{
		"desc": desc,
		"ok":   fmt.Sprint(itemOK),
	}

	// 行业均值水平估值
	checkItemName = "行业均值水平估值"
	itemOK = true
	desc = stock.JZPG.GetValuationScore()
	if stock.JZPG.GetValuationScore() == "高于行业均值水平" {
		ok = false
		itemOK = false
	}
	result[checkItemName] = map[string]string{
		"desc": desc,
		"ok":   fmt.Sprint(itemOK),
	}

	// 市盈率、市净率、市现率、市销率全部估值较高
	checkItemName = "四率估值"
	itemOK = true
	allHighValuation := true
	valuationDesc := []string{}
	for k, v := range stock.ValuationMap {
		valuationDesc = append(valuationDesc, k+v)
	}
	for _, v := range stock.ValuationMap {
		if v != "估值较高" {
			allHighValuation = false
			break
		}
	}
	if allHighValuation {
		ok = false
		itemOK = false
	}
	result[checkItemName] = map[string]string{
		"desc": strings.Join(valuationDesc, "\n"),
		"ok":   fmt.Sprint(itemOK),
	}

	// 股价低于合理价格
	if c.Options.IsCheckPriceByCalc {
		checkItemName = "合理股价"
		itemOK = true
		desc = fmt.Sprintf("最新股价:%f 合理价:%f", stock.GetPrice(), stock.RightPrice)
		if stock.RightPrice != -1 && stock.GetPrice() > stock.RightPrice {
			ok = false
			itemOK = false
		}
		result[checkItemName] = map[string]string{
			"desc": desc,
			"ok":   fmt.Sprint(itemOK),
		}
	}

	// 负债率低于 MaxDebtRatio （可选条件），金融股不检测该项
	if !goutils.IsStrInSlice(stock.GetOrgType(), []string{"银行", "保险"}) {
		if c.Options.MaxDebtAssetRatio != 0 && len(stock.HistoricalFinaMainData) > 0 {
			checkItemName = "负债率"
			itemOK = true
			fzl := stock.HistoricalFinaMainData[0].Zcfzl
			desc = fmt.Sprintf("负债率:%f", fzl)
			if fzl > c.Options.MaxDebtAssetRatio {
				desc = fmt.Sprintf("负债率:%f 高于:%f", fzl, c.Options.MaxDebtAssetRatio)
				ok = false
				itemOK = false
			}
			result[checkItemName] = map[string]string{
				"desc": desc,
				"ok":   fmt.Sprint(itemOK),
			}
		}
	}

	// 历史波动率 （可选条件）
	if c.Options.MaxHV != 0 {
		checkItemName = "历史波动率"
		itemOK = true
		desc = fmt.Sprintf("历史波动率:%f", stock.HistoricalVolatility)
		if stock.HistoricalVolatility > c.Options.MaxHV {
			desc = fmt.Sprintf("历史波动率:%f 高于:%f", stock.HistoricalVolatility, c.Options.MaxHV)
			ok = false
			itemOK = false
		}
		result[checkItemName] = map[string]string{
			"desc": desc,
			"ok":   fmt.Sprint(itemOK),
		}
	}

	// 市值
	checkItemName = "市值"
	itemOK = true
	sz := goutils.YiWanString(stock.BaseInfo.TotalMarketCap)
	desc = fmt.Sprintf("市值:%s", sz)
	if stock.BaseInfo.TotalMarketCap < c.Options.MinTotalMarketCap*100000000 {
		desc = fmt.Sprintf("市值:%s 低于:%f亿", sz, c.Options.MinTotalMarketCap)
		ok = false
		itemOK = false
	}
	result[checkItemName] = map[string]string{
		"desc": desc,
		"ok":   fmt.Sprint(itemOK),
	}

	// 银行股特殊检测
	if stock.GetOrgType() == "银行" && len(stock.HistoricalFinaMainData) > 0 {
		fmd := stock.HistoricalFinaMainData[0]
		checkItemName = "总资产收益率 (ROA)"
		itemOK = true
		desc = fmt.Sprintf("最新 ROA:%f", stock.BaseInfo.ROA)
		if stock.BaseInfo.ROA < c.Options.BankMinROA {
			desc = fmt.Sprintf("ROA:%f 低于:%f", stock.BaseInfo.ROA, c.Options.BankMinROA)
			ok = false
			itemOK = false
		}
		result[checkItemName] = map[string]string{
			"desc": desc,
			"ok":   fmt.Sprint(itemOK),
		}

		checkItemName = "资本充足率"
		itemOK = true
		desc = fmt.Sprintf("资本充足率:%f", fmd.Newcapitalader)
		if fmd.Newcapitalader < c.Options.BankMinZBCZL {
			desc = fmt.Sprintf("资本充足率:%f 低于:%f", fmd.Newcapitalader, c.Options.BankMinZBCZL)
			ok = false
			itemOK = false
		}
		result[checkItemName] = map[string]string{
			"desc": desc,
			"ok":   fmt.Sprint(itemOK),
		}

		if c.Options.BankMaxBLDKL != 0 {
			checkItemName = "不良贷款率"
			itemOK = true
			desc = fmt.Sprintf("不良贷款率:%f", fmd.Newcapitalader)
			if fmd.NonPerLoan > c.Options.BankMaxBLDKL {
				desc = fmt.Sprintf("不良贷款率:%f 高于:%f", fmd.Newcapitalader, c.Options.BankMinZBCZL)
				ok = false
				itemOK = false
			}
			result[checkItemName] = map[string]string{
				"desc": desc,
				"ok":   fmt.Sprint(itemOK),
			}
		}

		checkItemName = "不良贷款拨备覆盖率"
		itemOK = true
		desc = fmt.Sprintf("不良贷款拨备覆盖率:%f", fmd.Newcapitalader)
		if fmd.Bldkbbl < c.Options.BankMinBLDKBBFGL {
			desc = fmt.Sprintf("不良贷款拨备覆盖率:%f 低于:%f", fmd.Newcapitalader, c.Options.BankMinZBCZL)
			ok = false
			itemOK = false
		}
		result[checkItemName] = map[string]string{
			"desc": desc,
			"ok":   fmt.Sprint(itemOK),
		}
	}

	// 毛利率稳定性 （只检测非金融股）
	if c.Options.IsCheckMLLStability && !goutils.IsStrInSlice(stock.GetOrgType(), []string{"银行", "保险"}) {
		mllList := stock.HistoricalFinaMainData.ValueList(ctx, "MLL", c.Options.CheckYears)
		checkItemName = "毛利率稳定性"
		itemOK = true
		desc = fmt.Sprintf("%d 年内毛利率:%v", c.Options.CheckYears, mllList)
		if stock.HistoricalFinaMainData.IsStability(ctx, "MLL", c.Options.CheckYears) {
			desc = fmt.Sprintf("%d 年内稳定性较差:%v", c.Options.CheckYears, mllList)
			ok = false
			itemOK = false
		}
		result[checkItemName] = map[string]string{
			"desc": desc,
			"ok":   fmt.Sprint(itemOK),
		}
	}

	// 净利率稳定性
	if c.Options.IsCheckMLLStability {
		jllList := stock.HistoricalFinaMainData.ValueList(ctx, "JLL", c.Options.CheckYears)
		checkItemName = "净利率稳定性"
		itemOK = true
		desc = fmt.Sprintf("%d 年内净利率:%v", c.Options.CheckYears, jllList)
		if stock.HistoricalFinaMainData.IsStability(ctx, "JLL", c.Options.CheckYears) {
			desc = fmt.Sprintf("%d 年内稳定性较差:%v", c.Options.CheckYears, jllList)
			ok = false
			itemOK = false
		}
		result[checkItemName] = map[string]string{
			"desc": desc,
			"ok":   fmt.Sprint(itemOK),
		}
	}

	// PEG
	if c.Options.MaxPEG != 0 {
		checkItemName = "PEG"
		itemOK = true
		desc = fmt.Sprintf("PEG:%v", stock.PEG)
		if stock.PEG > c.Options.MaxPEG {
			desc = fmt.Sprintf("PEG:%v 高于:%v", stock.PEG, c.Options.MaxPEG)
			ok = false
			itemOK = false
		}
		result[checkItemName] = map[string]string{
			"desc": desc,
			"ok":   fmt.Sprint(itemOK),
		}
	}

	// 本业营收比
	if c.Options.MinBYYSRatio != 0 && c.Options.MaxBYYSRatio != 0 {
		checkItemName = "本业营收比"
		itemOK = true
		desc = fmt.Sprintf("当前本业营收比:%v", stock.BYYSRatio)
		if stock.BYYSRatio > c.Options.MaxBYYSRatio || stock.BYYSRatio < c.Options.MinBYYSRatio {
			desc = fmt.Sprintf("当前本业营收比:%v 超出范围:%v-%v", stock.BYYSRatio, c.Options.MinBYYSRatio, c.Options.MaxBYYSRatio)
			ok = false
			itemOK = false
		}
		result[checkItemName] = map[string]string{
			"desc": desc,
			"ok":   fmt.Sprint(itemOK),
		}
	}

	// 审计意见
	if stock.FinaReportOpinion != nil {
		checkItemName = "财报审计意见"
		itemOK = true
		opinion := stock.FinaReportOpinion.(string)
		desc = opinion
		if opinion != "标准无保留意见" {
			ok = false
			itemOK = false
		}
		result[checkItemName] = map[string]string{
			"desc": desc,
			"ok":   fmt.Sprint(itemOK),
		}
	}

	return
}
