// 按我的选股指标获取股票数据，对优质公司进行初步筛选（好公司不代表股价涨）
// 1. 净资产收益率， ROE_WEIGHT
// 2. 净利润增长率 ， NETPROFIT_YOY_RATIO
// 3. 营收增长率 ， TOI_YOY_RATIO
// 4. 最新股息率， ZXGXL
// 5. 净利润 3 年复合增长率， NETPROFIT_GROWTHRATE_3Y
// 6. 营收 3 年复合增长率， INCOME_GROWTHRATE_3Y
// 7. 预测净利润同比增长， PREDICT_NETPROFIT_RATIO
// 8. 预测营收同比增长， PREDICT_INCOME_RATIO
// 9. 上市以来年化收益率 ， LISTING_YIELD_YEAR
// 10. 总市值， TOTAL_MARKET_CAP
// 11. 是否按行业选择， INDUSTRY
// 12. 按股价（低股价 10-30 元)， NEW_PRICE
// 13. 上市时间是否大于 5 年，@LISTING_DATE="OVER5Y"
// 14. 市净率不小于 1 ， PBNEWMRQ

package eastmoney

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"go.uber.org/zap"
)

const (
	// YI 亿
	YI = float64(100000000)
	// WAN 万
	WAN = float64(10000)
)

// Filter 我的选股指标
type Filter struct {
	// ------ 最重要的指标！！！------
	// 最低净资产收益率（%）
	MinROE float64

	// ------ 必要参数 ------
	// 最低净利润增长率（%）
	MinNetprofitYoyRatio float64
	// 最低营收增长率（%）
	MinToiYoyRatio float64
	// 最低最新股息率（%）
	MinZXGXL float64
	// 最低净利润 3 年复合增长率（%）
	MinNetprofitGrowthrate3Y float64
	// 最低营收 3 年复合增长率（%）
	MinIncomeGrowthrate3Y float64
	// 最低预测净利润同比增长（%）
	MinPredictNetprofitRatio float64
	// 最低预测营收同比增长（%）
	MinPredictIncomeRatio float64
	// 最低上市以来年化收益率（%）
	MinListingYieldYear float64
	// 最低市净率
	MinPBNewMRQ float64

	// ------ 可选参数 ------
	// 最低总市值（亿）
	MinTotalMarketCap float64
	// 行业名（可选参数，不设置搜全行业）
	Industry string
	// 股价范围最小值（元）
	MinPrice float64
	// 股价范围最大值（元）
	MaxPrice float64
	// 上市时间是否超过 5 年
	ListingOver5Y bool
	// 是否排除创业板 300XXX
	ExcludeCYB bool
	// 是否排除科创板 688XXX
	ExcludeKCB bool
}

// String 转为字符串的请求参数
func (f Filter) String() string {
	filter := ""
	// 必要参数
	filter += fmt.Sprintf(`(ROE_WEIGHT>=%f)`, f.MinROE)
	filter += fmt.Sprintf(`(NETPROFIT_YOY_RATIO>=%f)`, f.MinNetprofitYoyRatio)
	filter += fmt.Sprintf(`(TOI_YOY_RATIO>=%f)`, f.MinToiYoyRatio)
	filter += fmt.Sprintf(`(ZXGXL>=%f)`, f.MinZXGXL)
	filter += fmt.Sprintf(`(NETPROFIT_GROWTHRATE_3Y>=%f)`, f.MinNetprofitGrowthrate3Y)
	filter += fmt.Sprintf(`(INCOME_GROWTHRATE_3Y>=%f)`, f.MinIncomeGrowthrate3Y)
	filter += fmt.Sprintf(`(PREDICT_NETPROFIT_RATIO>=%f)`, f.MinPredictNetprofitRatio)
	filter += fmt.Sprintf(`(PREDICT_INCOME_RATIO>=%f)`, f.MinPredictIncomeRatio)
	filter += fmt.Sprintf(`(LISTING_YIELD_YEAR>=%f)`, f.MinListingYieldYear)
	filter += fmt.Sprintf(`(PBNEWMRQ>=%f)`, f.MinPBNewMRQ)
	// 可选参数
	if f.MinTotalMarketCap != 0 {
		filter += fmt.Sprintf(`(TOTAL_MARKET_CAP>=%f)`, f.MinTotalMarketCap*100000000)
	}
	if f.Industry != "" {
		filter += fmt.Sprintf(`(INDUSTRY in ("%s"))`, f.Industry)
	}
	if f.MinPrice != 0 {
		filter += fmt.Sprintf(`(NEW_PRICE>=%f))`, f.MinPrice)
	}
	if f.MaxPrice != 0 {
		filter += fmt.Sprintf(`(NEW_PRICE<=%f))`, f.MaxPrice)
	}
	if f.ListingOver5Y {
		filter += `(@LISTING_DATE="OVER5Y")`
	}
	return filter
}

var (
	// DefaultFilter 默认指标值
	DefaultFilter = Filter{
		MinROE:                   8.0,
		MinNetprofitYoyRatio:     0.0,
		MinToiYoyRatio:           0.0,
		MinZXGXL:                 0.0,
		MinNetprofitGrowthrate3Y: 0.0,
		MinIncomeGrowthrate3Y:    0.0,
		MinPredictNetprofitRatio: 0.0,
		MinPredictIncomeRatio:    0.0,
		MinListingYieldYear:      0.0,
		MinPBNewMRQ:              0.0,
		MinTotalMarketCap:        0.0,
		Industry:                 "",
		MinPrice:                 0.0,
		MaxPrice:                 0.0,
		ListingOver5Y:            false,
		ExcludeCYB:               false,
		ExcludeKCB:               false,
	}
)

// StockInfo 接口返回的股票信息结构
type StockInfo struct {
	// 股票代码：带后缀
	Secucode string `json:"SECUCODE"`
	// 股票代码：无后缀
	SecurityCode string `json:"SECURITY_CODE"`
	// 股票名
	SecurityNameAbbr string `json:"SECURITY_NAME_ABBR"`
	// 行业
	Industry string `json:"INDUSTRY"`
	// 最新一期 ROE
	RoeWeight float64 `json:"ROE_WEIGHT"`
	// 净利润增长率（%）
	NetprofitYoyRatio float64 `json:"NETPROFIT_YOY_RATIO"`
	// 营收增长率（%）
	ToiYoyRatio float64 `json:"TOI_YOY_RATIO"`
	// 最新股息率
	Zxgxl float64 `json:"ZXGXL"`
	// 净利润 3 年复合增长率
	NetprofitGrowthrate3Y float64 `json:"NETPROFIT_GROWTHRATE_3Y"`
	//营收 3 年复合增长率
	IncomeGrowthrate3Y float64 `json:"INCOME_GROWTHRATE_3Y"`
	// 预测净利润同比增长
	PredictNetprofitRatio float64 `json:"PREDICT_NETPROFIT_RATIO"`
	// 预测营收同比增长
	PredictIncomeRatio float64 `json:"PREDICT_INCOME_RATIO"`
	// 上市以来年化收益率
	ListingYieldYear float64 `json:"LISTING_YIELD_YEAR"`
	// 总市值
	TotalMarketCap float64 `json:"TOTAL_MARKET_CAP"`
	// 最新价（元）
	NewPrice float64 `json:"NEW_PRICE"`
	// 市净率
	PBNewMRQ float64 `json:"PBNEWMRQ"`
}

// TotalMarketCapString 总市值可读字符串
func (s StockInfo) TotalMarketCapString() string {
	yi := s.TotalMarketCap / YI
	if yi >= 1 {
		return fmt.Sprintf("%.2f 亿", yi)
	}
	wan := s.TotalMarketCap / WAN
	if wan >= 1 {
		return fmt.Sprintf("%.2f 万", wan)
	}
	return fmt.Sprint(s.TotalMarketCap)
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

// QuerySelectedStocks 按选股指标默认值筛选股票
func (e EastMoney) QuerySelectedStocks(ctx context.Context) (StockInfoList, error) {
	return e.QuerySelectedStocksWithFilter(ctx, DefaultFilter)
}

// QuerySelectedStocksWithFilter 自定义选股指标值筛选股票
func (e EastMoney) QuerySelectedStocksWithFilter(ctx context.Context, filter Filter) (StockInfoList, error) {
	apiurl := "https://datacenter.eastmoney.com/stock/selection/api/data/get/"
	reqData := map[string]string{
		"source": "SELECT_SECURITIES",
		"client": "APP",
		"type":   "RPTA_APP_STOCKSELECT",
		"sty":    "SECUCODE,SECURITY_CODE,SECURITY_NAME_ABBR,INDUSTRY,ROE_WEIGHT,NETPROFIT_YOY_RATIO,TOI_YOY_RATIO,ZXGXL,NETPROFIT_GROWTHRATE_3Y,INCOME_GROWTHRATE_3Y,PREDICT_NETPROFIT_RATIO,PREDICT_INCOME_RATIO,LISTING_YIELD_YEAR,TOTAL_MARKET_CAP,NEW_PRICE,PBNEWMRQ",
		"filter": filter.String(),
		"p":      "1",      // page
		"ps":     "100000", // page size
	}
	logging.Debug(ctx, "EastMoney QuerySelectedStocksWithFilter "+apiurl+" begin", zap.Any("reqData", reqData))
	beginTime := time.Now()
	req, err := goutils.NewHTTPMultipartReq(ctx, apiurl, reqData)
	if err != nil {
		return nil, err
	}
	resp := RespSelectStocks{}
	if err := goutils.HTTPPOST(ctx, e.HTTPClient, req, &resp); err != nil {
		return nil, err
	}
	latency := time.Now().Sub(beginTime).Milliseconds()
	logging.Debug(
		ctx,
		"EastMoney SelectStocksWithFilter "+apiurl+" end",
		zap.Int64("latency(ms)", latency),
		zap.Any("resp", resp),
	)
	if resp.Code != 0 {
		return nil, fmt.Errorf("%#v", resp)
	}
	result := StockInfoList{}
	for _, i := range resp.Result.Data {
		// 排除创业板
		if filter.ExcludeCYB && strings.HasPrefix(i.Secucode, "300") {
			logging.Debugf(ctx, "EastMoney SelectStocksWithFilter ExcludeCYB %s %s", i.SecurityNameAbbr, i.Secucode)
			continue
		}
		// 排除科创板
		if filter.ExcludeKCB && strings.HasPrefix(i.Secucode, "688") {
			logging.Debugf(ctx, "EastMoney SelectStocksWithFilter ExcludeKCB %s %s", i.SecurityNameAbbr, i.Secucode)
			continue
		}
		result = append(result, i)
	}
	result.SortByROE()
	return result, nil
}
