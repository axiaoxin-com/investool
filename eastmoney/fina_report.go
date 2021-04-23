// 获取财报数据

package eastmoney

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"go.uber.org/zap"
)

// MainFinaData 财报主要指标
type MainFinaData struct {
	// 股票代码
	Secucode string `json:"SECUCODE"`
	// 股票代码
	SecurityCode string `json:"SECURITY_CODE"`
	// 股票名称
	SecurityNameAbbr string `json:"SECURITY_NAME_ABBR"`
	OrgCode          string `json:"ORG_CODE"`
	OrgType          string `json:"ORG_TYPE"`
	// 报告期
	ReportDate string `json:"REPORT_DATE"`
	// 财报类型：年报、三季报、中报、一季报
	ReportType string `json:"REPORT_TYPE"`
	// 财报名称: 2021一季报
	ReportDateName string `json:"REPORT_DATE_NAME"`
	// 财报年份：2021
	ReportYear       string `json:"REPORT_YEAR"`
	SecurityTypeCode string `json:"SECURITY_TYPE_CODE"`
	NoticeDate       string `json:"NOTICE_DATE"`
	// 财报更新时间
	UpdateDate string `json:"UPDATE_DATE"`
	// 货币类型：CNY
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

	// ------ ??? ------
	Totaldeposits     interface{} `json:"TOTALDEPOSITS"`
	Grossloans        interface{} `json:"GROSSLOANS"`
	Ltdrr             interface{} `json:"LTDRR"`
	Newcapitalader    interface{} `json:"NEWCAPITALADER"`
	Hxyjbczl          interface{} `json:"HXYJBCZL"`
	Nonperloan        interface{} `json:"NONPERLOAN"`
	Bldkbbl           interface{} `json:"BLDKBBL"`
	Nzbje             interface{} `json:"NZBJE"`
	TotalRoi          interface{} `json:"TOTAL_ROI"`
	NetRoi            interface{} `json:"NET_ROI"`
	EarnedPremium     interface{} `json:"EARNED_PREMIUM"`
	CompensateExpense interface{} `json:"COMPENSATE_EXPENSE"`
	SurrenderRateLife interface{} `json:"SURRENDER_RATE_LIFE"`
	SolvencyAr        interface{} `json:"SOLVENCY_AR"`
	Jzb               interface{} `json:"JZB"`
	Jzc               interface{} `json:"JZC"`
	Jzbjzc            interface{} `json:"JZBJZC"`
	Zygpgmjzc         interface{} `json:"ZYGPGMJZC"`
	Zygdsylzqjzb      interface{} `json:"ZYGDSYLZQJZB"`
	Yyfxzb            interface{} `json:"YYFXZB"`
	Jjywfxzb          interface{} `json:"JJYWFXZB"`
	Zqzyywfxzb        interface{} `json:"ZQZYYWFXZB"`
	Zqcxywfxzb        interface{} `json:"ZQCXYWFXZB"`
	Rzrqywfxzb        interface{} `json:"RZRQYWFXZB"`
	NbvLife           interface{} `json:"NBV_LIFE"`
	NbvRate           interface{} `json:"NBV_RATE"`
	NhjzCurrentAmt    interface{} `json:"NHJZ_CURRENT_AMT"`
}

// HistoryMainFinaData 主要指标历史数据列表
type HistoryMainFinaData []MainFinaData

// FilterByReportType 按财报类型过滤：一季报，中报，三季报，年报
func (h HistoryMainFinaData) FilterByReportType(ctx context.Context, reportType string) HistoryMainFinaData {
	result := HistoryMainFinaData{}
	for _, i := range h {
		if i.ReportType == reportType {
			result = append(result, i)
		}
	}
	return result
}

// FilterByReportYear 按财报年份过滤：2021
func (h HistoryMainFinaData) FilterByReportYear(ctx context.Context, reportYear string) HistoryMainFinaData {
	result := HistoryMainFinaData{}
	for _, i := range h {
		if i.ReportYear == reportYear {
			result = append(result, i)
		}
	}
	return result
}

// RespMainFinaData 接口返回 json 结构
type RespMainFinaData struct {
	Version string `json:"version"`
	Result  struct {
		Pages int                 `json:"pages"`
		Data  HistoryMainFinaData `json:"data"`
		Count int                 `json:"count"`
	} `json:"result"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// QueryMainFinaData 获取财报主要指标
func (e EastMoney) QueryMainFinaData(ctx context.Context, secuCode string) (HistoryMainFinaData, error) {
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
	logging.Debug(ctx, "EastMoney QueryMainFinaData "+apiurl+" begin", zap.Any("params", params))
	beginTime := time.Now()
	apiurl, err := goutils.NewHTTPGetURLWithQueryString(ctx, apiurl, params)
	if err != nil {
		return nil, err
	}
	resp := RespMainFinaData{}
	if err := goutils.HTTPGET(ctx, e.HTTPClient, apiurl, &resp); err != nil {
		return nil, err
	}
	latency := time.Now().Sub(beginTime).Milliseconds()
	logging.Debug(ctx, "EastMoney QueryMainFinaData "+apiurl+" end", zap.Int64("latency(ms)", latency), zap.Any("resp", resp))
	if resp.Code != 0 {
		return nil, fmt.Errorf("%#v", resp)
	}
	return resp.Result.Data, nil
}
