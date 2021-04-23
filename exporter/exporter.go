// Package exporter 导出各类型的数据结果
package exporter

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"sort"

	"github.com/axiaoxin-com/x-stock/model"
)

// Data 数据模板
type Data struct {
	// 股票名
	Name string `json:"name"`
	// 股票代码
	Code string `json:"code"`
	// 所属行业
	Industry string `json:"industry"`
	// 所属概念
	Concept string `json:"concept"`
	// 公司简介
	CompanyProfile string `json:"company_profile"`
	// 题材关键词
	Keywords string `json:"keywords"`
	// 主营构成
	MainForms string `json:"main_forms"`
	// 总市值（数字）
	TotalMarketCap float64 `json:"total_market_cap"`
	// 总市值（字符串）
	TotalMarketCapString string `json:"total_market_cap_string"`
	// 市盈率估值
	ValuationSYL string `json:"valuation_syl"`
	// 市净率估值
	ValuationSJL string `json:"valuation_sjl"`
	// 市销率估值
	ValuationSXOL string `json:"valuation_sxol"`
	// 市现率估值
	ValuationSXNL string `json:"valuation_sxnl"`
	// 综合估值状态
	ValuationStatusDesc string `json:"valuation_status_desc"`
	// 最新一期 ROE
	LatestROE float64 `json:"latest_roe"`
	// 模板当日价格
	Price float64 `json:"price"`
	// 估算合理价格
	RightPrice float64 `json:"right_price"`
	// 当前价格是否合理
	HasRightPrice bool `json:"has_right_price"`
	// 合理价格与现价的价格差
	PriceSpace float64 `json:"price_space"`
	// 历史波动率
	HV float64 `json:"hv"`
	// 净利润增长率（%）
	NetprofitYoyRatio float64 `json:"netprofit_yoy_ratio"`
	// 营收增长率（%）
	ToiYoyRatio float64 `json:"toi_yoy_ratio"`
	// 最新股息率（%）
	ZXGXL float64 `json:"zxgxl"`
	// 净利润 3 年复合增长率（%）
	NetprofitGrowthrate3Y float64 `json:"netprofit_growthrate_3_y"`
	// 营收 3 年复合增长率（%）
	IncomeGrowthrate3Y float64 `json:"income_growthrate_3_y"`
	// 预测净利润同比增长（%）
	PredictNetprofitRatio float64 `json:"predict_netprofit_ratio"`
	// 预测营收同比增长（%）
	PredictIncomeRatio float64 `json:"predict_income_ratio"`
	// 上市以来年化收益率（%）
	ListingYieldYear float64 `json:"listing_yield_year"`
	// 预约财报披露日期
	FinaAppointPublishDate string `json:"fina_appoint_publish_date"`
	// 机构评级
	OrgRating string `json:"org_rating"`
	// 每股收益预测
	EPSPredict string `json:"eps_predict"`
}

// NewData 创建 Data 对象
func NewData(stock model.Stock) Data {
	hasRightPrice := false
	if stock.RightPrice > 0 {
		hasRightPrice = true
	}
	return Data{
		Name:                   stock.BaseInfo.SecurityNameAbbr,
		Code:                   stock.BaseInfo.Secucode,
		Industry:               stock.BaseInfo.Industry,
		Concept:                stock.CompanyProfile.Concept,
		CompanyProfile:         stock.CompanyProfile.ProfileString(),
		Keywords:               stock.CompanyProfile.KeywordsString(),
		MainForms:              stock.CompanyProfile.MainFormsString(),
		TotalMarketCap:         stock.BaseInfo.TotalMarketCap,
		TotalMarketCapString:   stock.BaseInfo.TotalMarketCapString(),
		ValuationSYL:           stock.ValuationMap["市盈率"],
		ValuationSJL:           stock.ValuationMap["市净率"],
		ValuationSXOL:          stock.ValuationMap["市销率"],
		ValuationSXNL:          stock.ValuationMap["市现率"],
		ValuationStatusDesc:    stock.ValuationStatusDesc(),
		LatestROE:              stock.BaseInfo.RoeWeight,
		Price:                  stock.BaseInfo.NewPrice,
		RightPrice:             stock.RightPrice,
		HasRightPrice:          hasRightPrice,
		PriceSpace:             stock.RightPrice - stock.BaseInfo.NewPrice,
		HV:                     stock.HistoricalVolatility,
		NetprofitYoyRatio:      stock.BaseInfo.NetprofitYoyRatio,
		ToiYoyRatio:            stock.BaseInfo.ToiYoyRatio,
		ZXGXL:                  stock.BaseInfo.Zxgxl,
		NetprofitGrowthrate3Y:  stock.BaseInfo.NetprofitGrowthrate3Y,
		IncomeGrowthrate3Y:     stock.BaseInfo.IncomeGrowthrate3Y,
		PredictNetprofitRatio:  stock.BaseInfo.PredictNetprofitRatio,
		PredictIncomeRatio:     stock.BaseInfo.PredictIncomeRatio,
		ListingYieldYear:       stock.BaseInfo.ListingYieldYear,
		FinaAppointPublishDate: stock.FinaAppointPublishDate,
		OrgRating:              stock.OrgRatingList.String(),
		EPSPredict:             stock.ProfitPredictList.String(),
	}
}

// ExportData 要导出的数据列表
type ExportData []Data

// SortByROE 股票列表按 ROE 排序
func (e ExportData) SortByROE() {
	sort.Slice(e, func(i, j int) bool {
		return e[i].LatestROE > e[j].LatestROE
	})
}

// SortByPrice 股票列表按股价排序
func (e ExportData) SortByPrice() {
	sort.Slice(e, func(i, j int) bool {
		return e[i].Price < e[j].Price
	})
}

// SortByZXGXL 股票列表按最新股息率排序
func (e ExportData) SortByZXGXL() {
	sort.Slice(e, func(i, j int) bool {
		return e[i].ZXGXL > e[j].ZXGXL
	})
}

// SortByPriceSpace 股票列表按合理价格空间排序
func (e ExportData) SortByPriceSpace() {
	sort.Slice(e, func(i, j int) bool {
		return e[i].PriceSpace > e[j].PriceSpace
	})
}

// InitExportData 创建导出数据列表
func InitExportData(ctx context.Context, stocks model.StockList) (result ExportData) {
	for _, s := range stocks {
		result = append(result, NewData(s))
	}
	return
}

// ExportJSON 数据导出为JSON文件
func (e ExportData) ExportJSON(ctx context.Context, filename string) (result []byte, err error) {
	result, err = json.MarshalIndent(e, "", "  ")
	err = ioutil.WriteFile(filename, result, 0644)
	return
}
