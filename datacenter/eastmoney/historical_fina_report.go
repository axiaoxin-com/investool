// 获取财报数据

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

// FinaMainData 财报主要指标
type FinaMainData struct {
	// 股票代码: 111111.SZ
	Secucode string `json:"SECUCODE"`
	// 股票代码: 111111
	SecurityCode string `json:"SECURITY_CODE"`
	// 股票名称
	SecurityNameAbbr string `json:"SECURITY_NAME_ABBR"`
	OrgCode          string `json:"ORG_CODE"`
	OrgType          string `json:"ORG_TYPE"`
	// 报告期
	ReportDate string `json:"REPORT_DATE"`
	// 财报类型：年报、三季报、中报、一季报
	ReportType string `json:"REPORT_TYPE"`
	// 财报名称: 2021 一季报
	ReportDateName string `json:"REPORT_DATE_NAME"`
	// 财报年份： 2021
	ReportYear       string `json:"REPORT_YEAR"`
	SecurityTypeCode string `json:"SECURITY_TYPE_CODE"`
	NoticeDate       string `json:"NOTICE_DATE"`
	// 财报更新时间
	UpdateDate string `json:"UPDATE_DATE"`
	// 货币类型： CNY
	Currency string `json:"CURRENCY"`

	// ------ 每股指标 ------
	// 基本每股收益（元）
	Epsjb float64 `json:"EPSJB"`
	// 基本每股收益同比增长（%）
	Epsjbtz float64 `json:"EPSJBTZ"`
	// 扣非每股收益（元）
	Epskcjb float64 `json:"EPSKCJB"`
	// 稀释每股收益（元）
	Epsxs float64 `json:"EPSXS"`
	// 每股净资产（元）
	Bps float64 `json:"BPS"`
	// 每股净资产同比增长（%）
	Bpstz float64 `json:"BPSTZ"`
	// 每股资本公积（元）
	Mgzbgj float64 `json:"MGZBGJ"`
	// 每股资本公积同比增长（%）
	Mgzbgjtz float64 `json:"MGZBGJTZ"`
	// 每股未分配利润（元）
	Mgwfplr float64 `json:"MGWFPLR"`
	// 每股未分配利润同比增长（%）
	Mgwfplrtz float64 `json:"MGWFPLRTZ"`
	// 每股经营现金流（元）
	Mgjyxjje float64 `json:"MGJYXJJE"`
	// 每股经营现金流同比增长（%）
	Mgjyxjjetz float64 `json:"MGJYXJJETZ"`

	// ------ 成长能力指标 ------
	// 营业总收入（元）
	Totaloperatereve float64 `json:"TOTALOPERATEREVE"`
	// 营业总收入同比增长（%）
	Totaloperaterevetz float64 `json:"TOTALOPERATEREVETZ"`
	// 毛利润（元）
	Mlr float64 `json:"MLR"`
	// 归属净利润（元）
	Parentnetprofit float64 `json:"PARENTNETPROFIT"`
	// 归属净利润同比增长（%）
	Parentnetprofittz float64 `json:"PARENTNETPROFITTZ"`
	// 扣非净利润（元）
	Kcfjcxsyjlr float64 `json:"KCFJCXSYJLR"`
	// 扣非净利润同比增长（%）
	Kcfjcxsyjlrtz float64 `json:"KCFJCXSYJLRTZ"`
	// 营业总收入滚动环比增长（%）
	Yyzsrgdhbzc float64 `json:"YYZSRGDHBZC"`
	// 归属净利润滚动环比增长（%）
	Netprofitrphbzc float64 `json:"NETPROFITRPHBZC"`
	// 扣非净利润滚动环比增长（%）
	Kfjlrgdhbzc float64 `json:"KFJLRGDHBZC"`

	// ------ 盈利能力指标 ------
	// 净资产收益率（加权）（%）
	Roejq float64 `json:"ROEJQ"`
	// 净资产收益率（扣非/加权）（%）
	Roekcjq float64 `json:"ROEKCJQ"`
	// 净资产收益率同比增长（%）
	Roejqtz float64 `json:"ROEJQTZ"`
	// 总资产收益率（加权）（%）
	Zzcjll float64 `json:"ZZCJLL"`
	// 总资产收益率同比增长（%）
	Zzcjlltz float64 `json:"ZZCJLLTZ"`
	// 投入资本回报率（%）
	Roic float64 `json:"ROIC"`
	// 投入资本回报率同比增长（%）
	Roictz float64 `json:"ROICTZ"`
	// 毛利率（%）
	Xsmll float64 `json:"XSMLL"`
	// 净利率（%）
	Xsjll float64 `json:"XSJLL"`

	// ------ 收益质量指标 ------
	// 预收账款/营业收入
	Yszkyysr interface{} `json:"YSZKYYSR"`
	// 销售净现金流/营业收入
	Xsjxlyysr float64 `json:"XSJXLYYSR"`
	// 经营净现金流/营业收入
	Jyxjlyysr float64 `json:"JYXJLYYSR"`
	// 实际税率（%）
	Taxrate float64 `json:"TAXRATE"`

	// ------ 财务风险指标 ------
	// 流动比率
	Ld float64 `json:"LD"`
	// 速动比率
	Sd float64 `json:"SD"`
	// 现金流量比率
	Xjllb float64 `json:"XJLLB"`
	// 资产负债率（%）
	Zcfzl float64 `json:"ZCFZL"`
	// 资产负债率同比增长（%）
	Zcfzltz float64 `json:"ZCFZLTZ"`
	// 权益乘数
	Qycs float64 `json:"QYCS"`
	// 产权比率
	Cqbl float64 `json:"CQBL"`

	// ------ 营运能力指标 ------
	// 总资产周转天数（天）
	Zzczzts float64 `json:"ZZCZZTS"`
	// 存货周转天数（天）
	Chzzts float64 `json:"CHZZTS"`
	// 应收账款周转天数（天）
	Yszkzzts float64 `json:"YSZKZZTS"`
	// 总资产周转率（次）
	Toazzl float64 `json:"TOAZZL"`
	// 存货周转率（次）
	Chzzl float64 `json:"CHZZL"`
	// 应收账款周转率（次）
	Yszkzzl float64 `json:"YSZKZZL"`
}

// HistoricalFinaMainData 主要指标历史数据列表
type HistoricalFinaMainData []FinaMainData

// FilterByReportType 按财报类型过滤：一季报，中报，三季报，年报
func (h HistoricalFinaMainData) FilterByReportType(ctx context.Context, reportType string) HistoricalFinaMainData {
	result := HistoricalFinaMainData{}
	for _, i := range h {
		if i.ReportType == reportType {
			result = append(result, i)
		}
	}
	return result
}

// FilterByReportYear 按财报年份过滤： 2021
func (h HistoricalFinaMainData) FilterByReportYear(ctx context.Context, reportYear int) HistoricalFinaMainData {
	result := HistoricalFinaMainData{}
	year := fmt.Sprint(reportYear)
	for _, i := range h {
		if i.ReportYear == year {
			result = append(result, i)
		}
	}
	return result
}

// ROEList 获取历史 roe，最新的在最前面
func (h HistoricalFinaMainData) ROEList(ctx context.Context, count int) []float64 {
	r := []float64{}
	data := h.FilterByReportType(ctx, "年报")
	if len(data) == 0 {
		return r
	}
	if count > 0 {
		data = data[:count]
	}
	for _, i := range data {
		r = append(r, i.Roejq)
	}
	return r
}

// EPSList 获取历史 eps，最新的在最前面
func (h HistoricalFinaMainData) EPSList(ctx context.Context, count int) []float64 {
	r := []float64{}
	data := h.FilterByReportType(ctx, "年报")
	if len(data) == 0 {
		return r
	}
	if count > 0 {
		data = data[:count]
	}
	for _, i := range data {
		r = append(r, i.Epsjb)
	}
	return r
}

// RevenueList 获取历史营收，最新的在最前面
func (h HistoricalFinaMainData) RevenueList(ctx context.Context, count int) []float64 {
	r := []float64{}
	data := h.FilterByReportType(ctx, "年报")
	if len(data) == 0 {
		return r
	}
	if count > 0 {
		data = data[:count]
	}
	for _, i := range data {
		r = append(r, i.Totaloperatereve)
	}
	return r
}

// ProfitList 获取历史利润，最新的在最前面
func (h HistoricalFinaMainData) ProfitList(ctx context.Context, count int) []float64 {
	r := []float64{}
	data := h.FilterByReportType(ctx, "年报")
	if len(data) == 0 {
		return r
	}
	if count > 0 {
		data = data[:count]
	}
	for _, i := range data {
		r = append(r, i.Parentnetprofit)
	}
	return r
}

// IsIncreasingByYears roe/eps/revenue/profit 是否逐年递增
func (h HistoricalFinaMainData) IsIncreasingByYears(ctx context.Context, dataType string, years int) bool {
	data := h.FilterByReportType(ctx, "年报")
	dataLen := len(data)
	if dataLen == 0 {
		return false
	}
	if years > dataLen {
		years = dataLen
	}

	dataType = strings.ToUpper(dataType)
	increasing := true
	for i := 0; i < years-1; i++ {
		switch dataType {
		case "ROE":
			if data[i].Roejq <= data[i+1].Roejq {
				increasing = false
				break
			}
		case "EPS":
			if data[i].Epsjb <= data[i+1].Epsjb {
				increasing = false
				break
			}
		case "REVENUE":
			if data[i].Totaloperatereve <= data[i+1].Totaloperatereve {
				increasing = false
				break
			}
		case "PROFIT":
			if data[i].Parentnetprofit <= data[i+1].Parentnetprofit {
				increasing = false
				break
			}
		}
	}
	return increasing
}

// MidValue 历史年报 roe/eps 中位数
func (h HistoricalFinaMainData) MidValue(ctx context.Context, dataType string, years int) float64 {
	dataType = strings.ToUpper(dataType)
	values := []float64{}
	data := h.FilterByReportType(ctx, "年报")
	if years > 0 {
		data = data[:years]
	}
	dataLen := len(data)
	if dataLen == 0 {
		return 0
	}
	for _, i := range data {
		switch dataType {
		case "ROE":
			values = append(values, i.Roejq)
		case "EPS":
			values = append(values, i.Epsjb)
		}
	}
	sort.Float64s(values)
	mid := dataLen / 2
	if dataLen%2 == 0 {
		return (values[mid-1] + values[mid]) / 2
	}
	return values[mid]
}

// Q1RevenueIncreasingRatio 获取今年一季报的营收增长比 (%)
func (h HistoricalFinaMainData) Q1RevenueIncreasingRatio(ctx context.Context) (float64, error) {
	year := time.Now().Year()
	data := h.FilterByReportYear(ctx, year)
	if len(data) > 0 {
		return data[0].Totaloperaterevetz, nil
	}
	return 0, fmt.Errorf("%dQ1 report has not yet been published", year)
}

// RespFinaMainData 接口返回 json 结构
type RespFinaMainData struct {
	Version string `json:"version"`
	Result  struct {
		Pages int                    `json:"pages"`
		Data  HistoricalFinaMainData `json:"data"`
		Count int                    `json:"count"`
	} `json:"result"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// QueryHistoricalFinaMainData 获取财报主要指标，最新的在最前面
func (e EastMoney) QueryHistoricalFinaMainData(ctx context.Context, secuCode string) (HistoricalFinaMainData, error) {
	apiurl := "https://datacenter.eastmoney.com/securities/api/data/get"
	params := map[string]string{
		"filter": fmt.Sprintf(`(SECUCODE="%s")`, strings.ToUpper(secuCode)),
		"client": "APP",
		"source": "HSF10",
		"type":   "RPT_F10_FINANCE_MAINFINADATA",
		"sty":    "APP_F10_MAINFINADATA",
		"st":     "REPORT_DATE",
		"ps":     "100",
		"sr":     "-1",
	}
	logging.Debug(ctx, "EastMoney QueryHistoricalFinaMainData "+apiurl+" begin", zap.Any("params", params))
	beginTime := time.Now()
	apiurl, err := goutils.NewHTTPGetURLWithQueryString(ctx, apiurl, params)
	if err != nil {
		return nil, err
	}
	resp := RespFinaMainData{}
	err = goutils.HTTPGET(ctx, e.HTTPClient, apiurl, &resp)
	latency := time.Now().Sub(beginTime).Milliseconds()
	logging.Debug(
		ctx,
		"EastMoney QueryHistoricalFinaMainData "+apiurl+" end",
		zap.Int64("latency(ms)", latency),
		zap.Any("resp", resp),
	)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("%s %#v", secuCode, resp)
	}
	return resp.Result.Data, nil
}

// RespFinaPublishDate 财报披露日期接口返回结构
type RespFinaPublishDate struct {
	Version string `json:"version"`
	Result  struct {
		Pages int `json:"pages"`
		Data  []struct {
			SecurityCode       string      `json:"SECURITY_CODE"`
			SecurityNameAbbr   string      `json:"SECURITY_NAME_ABBR"`
			AppointPublishDate string      `json:"APPOINT_PUBLISH_DATE"`
			ReportDate         string      `json:"REPORT_DATE"`
			ActualPublishDate  interface{} `json:"ACTUAL_PUBLISH_DATE"`
			ReportTypeName     string      `json:"REPORT_TYPE_NAME"`
			IsPublish          string      `json:"IS_PUBLISH"`
		} `json:"data"`
		Count int `json:"count"`
	} `json:"result"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// QueryAppointFinaPublishDate 查询最新财报预约披露日期
func (e EastMoney) QueryAppointFinaPublishDate(ctx context.Context, securityCode string) (string, error) {
	apiurl := "https://datacenter.eastmoney.com/api/data/get"
	params := map[string]string{
		"filter": fmt.Sprintf(`(SECURITY_CODE="%s")`, strings.ToUpper(securityCode)),
		"client": "APP",
		"source": "DataCenter",
		"type":   "RPT_PUBLIC_BS_APPOIN",
		"sty":    "SECURITY_CODE,SECURITY_NAME_ABBR,APPOINT_PUBLISH_DATE,REPORT_DATE,ACTUAL_PUBLISH_DATE,REPORT_TYPE_NAME,IS_PUBLISH",
		"st":     "SECURITY_CODE,EITIME",
		"ps":     "20",
		"p":      "1",
		"sr":     "-1,-1",
	}
	logging.Debug(ctx, "EastMoney QueryAppointFinaPublishDate "+apiurl+" begin", zap.Any("params", params))
	beginTime := time.Now()
	apiurl, err := goutils.NewHTTPGetURLWithQueryString(ctx, apiurl, params)
	if err != nil {
		return "", err
	}
	resp := RespFinaPublishDate{}
	err = goutils.HTTPGET(ctx, e.HTTPClient, apiurl, &resp)
	latency := time.Now().Sub(beginTime).Milliseconds()
	logging.Debug(
		ctx,
		"EastMoney QueryAppointFinaPublishDate "+apiurl+" end",
		zap.Int64("latency(ms)", latency),
		zap.Any("resp", resp),
	)
	if err != nil {
		return "", err
	}
	if resp.Code != 0 {
		return "", fmt.Errorf("%s %#v", securityCode, resp)
	}
	if len(resp.Result.Data) > 0 {
		return resp.Result.Data[0].AppointPublishDate, nil
	}
	return "", nil
}
