// 按我的选股指标获取股票数据
// 1. 总市值 > 500 亿
// 2. 净资产收益率 >= 8%
// 3. 最新股息率 >= 0.000001%
// 4. 净利润 3 年复合增长率 > 0
// 5. 营收 3 年复合增长率 > 0
// 6. 预测净利润同比增长 > 0
// 7. 预测营收同比增长 > 0
// 8. 每股股利（税前） >= 0.000001
// 9. 上市以来年化收益率 > 0
//
// (INDUSTRY in ("行业名"))(TOTAL_MARKET_CAP>50000000000)(ROE_WEIGHT>=8)(ZXGXL>0.000001)(NETPROFIT_GROWTHRATE_3Y>0)(INCOME_GROWTHRATE_3Y>0)(PREDICT_NETPROFIT_RATIO>0)(PREDICT_INCOME_RATIO>0)(PAR_DIVIDEND_PRETAX>0.000001)(LISTING_YIELD_YEAR>0)

package eastmoney

import (
	"context"
	"fmt"
	"sort"
)

// Filter 我的选股指标
type Filter struct {
	// 总市值，单位：亿
	TotalMarketCap float64 `json:"total_market_cap"`
	// 净资产收益率，单位：%
	ROE float64 `json:"roe"`
	// 最新股息率，单位：%
	ZXGXL float64 `json:"zxgxl"`
	// 净利润 3 年复合增长率，单位：%
	NetprofitGrowthrate3Y float64 `json:"netprofit_growthrate_3_y"`
	// 营收 3 年复合增长率，单位：%
	IncomeGrowthrate3Y float64 `json:"income_growthrate_3_y"`
	// 预测净利润同比增长，单位：%
	PredictNetprofitRatio float64 `json:"predict_netprofit_ratio"`
	// 预测营收同比增长，单位：%
	PredictIncomeRatio float64 `json:"predict_income_ratio"`
	// 每股股利（税前），单位：元
	ParDividendPretax float64 `json:"par_dividend_pretax"`
	// 上市以来年化收益率，单位：%
	ListingYieldYear float64 `json:"listing_yield_year"`
	// 行业名
	Industry string `json:"industry"`
}

// ToString 转为字符串的请求参数
func (f Filter) ToString(ctx context.Context) string {
	filter := ""
	if f.Industry != "" {
		filter += fmt.Sprintf(`(INDUSTRY in ("%s"))`, f.Industry)
	}
	filter += fmt.Sprintf(`(TOTAL_MARKET_CAP>%f)`, f.TotalMarketCap*100000000)
	filter += fmt.Sprintf(`(ROE_WEIGHT>=%f)`, f.ROE)
	filter += fmt.Sprintf(`(ZXGXL>=%f)`, f.ZXGXL)
	filter += fmt.Sprintf(`(NETPROFIT_GROWTHRATE_3Y>%f)`, f.NetprofitGrowthrate3Y)
	filter += fmt.Sprintf(`(INCOME_GROWTHRATE_3Y>%f)`, f.IncomeGrowthrate3Y)
	filter += fmt.Sprintf(`(PREDICT_NETPROFIT_RATIO>%f)`, f.PredictNetprofitRatio)
	filter += fmt.Sprintf(`(PREDICT_INCOME_RATIO>%f)`, f.PredictIncomeRatio)
	filter += fmt.Sprintf(`(PAR_DIVIDEND_PRETAX>=%f)`, f.ParDividendPretax)
	filter += fmt.Sprintf(`(LISTING_YIELD_YEAR>%f)`, f.ListingYieldYear)
	return filter
}

// StockInfo 接口返回的股票信息结构
type StockInfo struct {
	// 股票代码
	Secucode string `json:"SECUCODE"`
	// 股票代码
	SecurityCode string `json:"SECURITY_CODE"`
	// 股票名
	SecurityNameAbbr string `json:"SECURITY_NAME_ABBR"`
	// 最新价（元）
	NewPrice float64 `json:"NEW_PRICE"`
	// 涨跌幅（%）
	ChangeRate float64 `json:"CHANGE_RATE"`
	// 行业
	Industry string `json:"INDUSTRY"`
	// 总市值
	TotalMarketCap int64 `json:"TOTAL_MARKET_CAP"`
	// 最新一期 ROE
	RoeWeight float64 `json:"ROE_WEIGHT"`
	// 最新股息率
	Zxgxl float64 `json:"ZXGXL"`
	// 净利润 3 年复合增长率
	NetprofitGrowthrate3Y float64 `json:"NETPROFIT_GROWTHRATE_3Y"`
	//营收 3 年复合增长率
	IncomeGrowthrate3Y float64 `json:"INCOME_GROWTHRATE_3Y"`
	// 预测净利润同比增长
	PredictNetprofitRatio float64 `json:"PREDICT_NETPROFIT_RATIO"`
	//	预测营收同比增长
	PredictIncomeRatio float64 `json:"PREDICT_INCOME_RATIO"`
	//每股股利（税前）
	ParDividendPretax float64 `json:"PAR_DIVIDEND_PRETAX"`
	// 上市以来年化收益率
	ListingYieldYear float64 `json:"LISTING_YIELD_YEAR"`
	// 最近交易日期
	MaxTradeDate string `json:"MAX_TRADE_DATE"`
}

// StockInfoList 股票列表
type StockInfoList []StockInfo

// SortByROE 股票列表按 ROE 排序
func (s StockInfoList) SortByROE() {
	sort.Slice(s, func(i, j int) bool {
		return s[i].RoeWeight > s[j].RoeWeight
	})
}

// RespSelectStocks 接口返回 json 结构
type RespSelectStocks struct {
	Result struct {
		Nextpage    bool          `json:"nextpage"`
		Currentpage int           `json:"currentpage"`
		Data        StockInfoList `json:"data"`
		Config      []struct {
			IndicatorName string `json:"INDICATOR_NAME"`
			Datatype      string `json:"DATATYPE"`
		} `json:"config"`
	} `json:"result"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

var (
	// DefaultFilter 默认指标值
	DefaultFilter = Filter{
		TotalMarketCap:        500.0,
		ROE:                   8.0,
		ZXGXL:                 0.000001,
		NetprofitGrowthrate3Y: 0.0,
		IncomeGrowthrate3Y:    0.0,
		PredictNetprofitRatio: 0.0,
		PredictIncomeRatio:    0.0,
		ParDividendPretax:     0.000001,
		ListingYieldYear:      0.0,
		Industry:              "",
	}
)

// SelectStocks 按选股指标默认值筛选股票
func (e EastMoney) SelectStocks(ctx context.Context) (StockInfoList, error) {
	return e.SelectStocksWithFilter(ctx, DefaultFilter)
}

// SelectStocksWithFilter 自定义选股指标值筛选股票
func (e EastMoney) SelectStocksWithFilter(ctx context.Context, filter Filter) (StockInfoList, error) {
	url := "https://datacenter.eastmoney.com/stock/selection/api/data/get/"
	reqData := map[string]string{
		"source": "SELECT_SECURITIES",
		"client": "APP",
		"type":   "RPTA_APP_STOCKSELECT",
		"sty":    "SECUCODE,SECURITY_CODE,SECURITY_NAME_ABBR,NEW_PRICE,CHANGE_RATE,INDUSTRY,TOTAL_MARKET_CAP,ROE_WEIGHT,ZXGXL,NETPROFIT_GROWTHRATE_3Y,INCOME_GROWTHRATE_3Y,PREDICT_NETPROFIT_RATIO,PREDICT_INCOME_RATIO,PAR_DIVIDEND_PRETAX,LISTING_YIELD_YEAR",
		"filter": filter.ToString(ctx),
		"p":      "1",      // page
		"ps":     "100000", // page size
	}
	resp := RespSelectStocks{}
	if err := e.Post(ctx, url, reqData, &resp); err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("%#v", resp)
	}
	result := StockInfoList{}
	for _, i := range resp.Result.Data {
		result = append(result, i)
	}
	return result, nil
}
