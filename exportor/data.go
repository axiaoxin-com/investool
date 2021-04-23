// 封装需要导出数据结果

package exportor

import (
	"context"
	"sort"
	"strings"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/x-stock/model"
)

// Data 数据模板
type Data struct {
	// 股票名
	Name string `json:"name"                      csv:"股票名"`
	// 股票代码
	Code string `json:"code"                      csv:"股票代码"`
	// 所属行业
	Industry string `json:"industry"                  csv:"所属行业"`
	// 题材关键词
	Keywords string `json:"keywords"                  csv:"题材关键词"`
	// 公司信息
	CompanyProfile string `json:"company_profile"           csv:"公司信息"`
	// 主营构成
	MainForms string `json:"main_forms"                csv:"主营构成"`
	// 总市值（数字）
	TotalMarketCap float64 `json:"total_market_cap"          csv:"-"`
	// 总市值（字符串）
	TotalMarketCapString string `json:"total_market_cap_string"   csv:"总市值"`
	// 市盈率估值
	ValuationSYL string `json:"valuation_syl"             csv:"市盈率估值"`
	// 市净率估值
	ValuationSJL string `json:"valuation_sjl"             csv:"市净率估值"`
	// 市销率估值
	ValuationSXOL string `json:"valuation_sxol"            csv:" 市销率估值"`
	// 市现率估值
	ValuationSXNL string `json:"valuation_sxnl"            csv:"市现率估值"`
	// 综合估值状态
	ValuationStatusDesc string `json:"valuation_status_desc"     csv:"综合估值状态"`
	// 最新一期 ROE
	LatestROE float64 `json:"latest_roe"                csv:"最新一期 ROE"`
	// 当时价格
	Price float64 `json:"price"                     csv:"价格"`
	// 估算合理价格
	RightPrice interface{} `json:"right_price"               csv:"估算合理价格"`
	// 当前价格是否合理
	HasRightPrice interface{} `json:"has_right_price"           csv:"当前价格是否合理"`
	// 合理价格与当时价的价格差
	PriceSpace interface{} `json:"price_space"               csv:"合理价格与当时价的价格差"`
	// 预约财报披露日期
	FinaAppointPublishDate string `json:"fina_appoint_publish_date" csv:"预约财报披露日期"`
	// 历史波动率 (%)
	HV float64 `json:"hv"                        csv:"历史波动率 (%)"`
	// 净利润增长率 (%)
	NetprofitYoyRatio float64 `json:"netprofit_yoy_ratio"       csv:"净利润增长率 (%)"`
	// 营收增长率 (%)
	ToiYoyRatio float64 `json:"toi_yoy_ratio"             csv:"营收增长率 (%)"`
	// 最新负债率 (%)
	ZXFZL float64 `json:"zxfzl"                     csv:"最新负债率 (%)"`
	// 最新股息率 (%)
	ZXGXL float64 `json:"zxgxl"                     csv:"最新股息率 (%)"`
	// 净利润 3 年复合增长率 (%)
	NetprofitGrowthrate3Y float64 `json:"netprofit_growthrate_3_y"  csv:"净利润 3 年复合增长率 (%)"`
	// 营收 3 年复合增长率 (%)
	IncomeGrowthrate3Y float64 `json:"income_growthrate_3_y"     csv:"营收 3 年复合增长率 (%)"`
	// 上市以来年化收益率 (%)
	ListingYieldYear float64 `json:"listing_yield_year"        csv:"上市以来年化收益率 (%)"`
	// 上市以来年化波动率 (%)
	ListingVolatilityYear float64 `json:"listing_volatility_year"   csv:"年化波动率 (%)"`
	// 预测净利润同比增长 (%)
	PredictNetprofitRatio float64 `json:"predict_netprofit_ratio"   csv:"预测净利润同比增长 (%)"`
	// 预测营收同比增长 (%)
	PredictIncomeRatio float64 `json:"predict_income_ratio"      csv:"预测营收同比增长 (%)"`
	// 机构评级
	OrgRating string `json:"org_rating"                csv:"机构评级"`
	// 每股收益预测
	EPSPredict string `json:"eps_predict"               csv:"每股收益预测"`
	// 上市时间
	ListingDate string `json:"listing_date"              csv:"上市时间"`
}

// GetHeaderValueMap 获取以 csv tag 为 key 的 Data map
func (d Data) GetHeaderValueMap() map[string]interface{} {
	return goutils.StructToMap(&d, "csv")
}

// GetHeaders 获取 csv tag 列表
func (d Data) GetHeaders() []string {
	return goutils.StructTagList(&d, "csv")
}

// NewData 创建 Data 对象
func NewData(stock model.Stock) Data {
	var rightPrice interface{} = "--"
	var hasRightPrice interface{} = "--"
	var priceSpace interface{} = "--"
	if stock.RightPrice > 0 {
		rightPrice = stock.RightPrice
		hasRightPrice = false
		if stock.GetPrice() < stock.RightPrice {
			hasRightPrice = true
		}
		priceSpace = stock.RightPrice - stock.GetPrice()
	}
	return Data{
		Name:                   stock.BaseInfo.SecurityNameAbbr,
		Code:                   stock.BaseInfo.Secucode,
		Industry:               stock.BaseInfo.Industry,
		Keywords:               stock.CompanyProfile.KeywordsString(),
		CompanyProfile:         stock.CompanyProfile.ProfileString(),
		MainForms:              stock.CompanyProfile.MainFormsString(),
		TotalMarketCap:         stock.BaseInfo.TotalMarketCap,
		TotalMarketCapString:   stock.BaseInfo.TotalMarketCapString(),
		ValuationSYL:           stock.ValuationMap["市盈率"],
		ValuationSJL:           stock.ValuationMap["市净率"],
		ValuationSXOL:          stock.ValuationMap["市销率"],
		ValuationSXNL:          stock.ValuationMap["市现率"],
		ValuationStatusDesc:    stock.ValuationStatusDesc(),
		LatestROE:              stock.BaseInfo.RoeWeight,
		Price:                  stock.GetPrice(),
		RightPrice:             rightPrice,
		HasRightPrice:          hasRightPrice,
		PriceSpace:             priceSpace,
		FinaAppointPublishDate: strings.Fields(stock.FinaAppointPublishDate)[0],
		HV:                     stock.HistoricalVolatility,
		ListingVolatilityYear:  stock.BaseInfo.ListingVolatilityYear,
		NetprofitYoyRatio:      stock.BaseInfo.NetprofitYoyRatio,
		ToiYoyRatio:            stock.BaseInfo.ToiYoyRatio,
		ZXFZL:                  stock.HistoricalFinaMainData[0].Zcfzl,
		ZXGXL:                  stock.BaseInfo.Zxgxl,
		NetprofitGrowthrate3Y:  stock.BaseInfo.NetprofitGrowthrate3Y,
		IncomeGrowthrate3Y:     stock.BaseInfo.IncomeGrowthrate3Y,
		ListingYieldYear:       stock.BaseInfo.ListingYieldYear,
		PredictNetprofitRatio:  stock.BaseInfo.PredictNetprofitRatio,
		PredictIncomeRatio:     stock.BaseInfo.PredictIncomeRatio,
		OrgRating:              stock.OrgRatingList.String(),
		EPSPredict:             stock.ProfitPredictList.String(),
		ListingDate:            stock.BaseInfo.ListingDate,
	}
}

// DataList 要导出的数据列表
type DataList []Data

// SortByROE 股票列表按 ROE 排序
func (d DataList) SortByROE() {
	sort.Slice(d, func(i, j int) bool {
		return d[i].LatestROE > d[j].LatestROE
	})
}

// SortByPrice 股票列表按股价排序
func (d DataList) SortByPrice() {
	sort.Slice(d, func(i, j int) bool {
		return d[i].Price < d[j].Price
	})
}

// SortByZXGXL 股票列表按最新股息率排序
func (d DataList) SortByZXGXL() {
	sort.Slice(d, func(i, j int) bool {
		return d[i].ZXGXL > d[j].ZXGXL
	})
}

// SortByHV 股票列表按历史波动率排序
func (d DataList) SortByHV() {
	sort.Slice(d, func(i, j int) bool {
		return d[i].HV > d[j].HV
	})
}

// IndustryList 获取行业分类列表
func (d DataList) GetIndustryList() []string {
	result := []string{}
	for _, stock := range d {
		if !goutils.IsStrInSlice(stock.Industry, result) {
			result = append(result, stock.Industry)
		}
	}
	return result
}

// NewDataList 创建要导出的数据列表
func NewDataList(ctx context.Context, stocks model.StockList) (result DataList) {
	for _, s := range stocks {
		result = append(result, NewData(s))
	}
	return
}
