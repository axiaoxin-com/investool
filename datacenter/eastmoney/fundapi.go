// 天天基金 API 封装

package eastmoney

import (
	"context"
	"fmt"
	"time"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"go.uber.org/zap"
)

// FundType 基金类型
type FundType int

const (
	// 基金类型定义

	// FundTypeALL 全部
	FundTypeALL FundType = 0
	// FundTypeStock 股票型
	FundTypeStock FundType = 25
	// FundTypeMix 混合型
	FundTypeMix FundType = 27
	// FundTypeIndex 指数型
	FundTypeIndex FundType = 26
	// FundTypeBond 债券型
	FundTypeBond FundType = 31
	// FundTypeCurrency 货币型
	FundTypeCurrency FundType = 35
	// FundTypeFOF FOF型
	FundTypeFOF FundType = 15
	// FundTypeQDII QDII型
	FundTypeQDII FundType = 6
	// FundTypeETF ETF
	FundTypeETF FundType = 3
	// FundTypeETFLink ETF联接
	FundTypeETFLink FundType = 33
	// FundTypeLOF LOF
	FundTypeLOF FundType = 4
	// FundTypeDealMoney 理财
	FundTypeDealMoney FundType = 2949
)

// FundNetListItem RespFundNetList Datas 列表元素
type FundNetListItem struct {
	// 基金代码
	Fcode string `json:"FCODE"`
	// 基金名
	Shortname string `json:"SHORTNAME"`
	// 001:股票型、指数型、ETF、ETF联接、LOF
	// 002:混合型、FOF
	// 003:债券型
	// 005:货币型
	// 007:QDII型
	Fundtype  string `json:"FUNDTYPE"`
	Bfundtype string `json:"BFUNDTYPE"`
	Feature   string `json:"FEATURE"`
	Fsrq      string `json:"FSRQ"`
	Gpsj      string `json:"GPSJ"`
	// 净值
	Zjl         string `json:"ZJL"`
	Opendate    string `json:"OPENDATE"`
	Targetyield string `json:"TARGETYIELD"`
	Dwjz        string `json:"DWJZ"`
	Hldwjz      string `json:"HLDWJZ"`
	Ljjz        string `json:"LJJZ"`
	// 日涨幅 (%)
	Rzdf    string `json:"RZDF"`
	Cycle1  string `json:"CYCLE1"`
	Isbuy   string `json:"ISBUY"`
	Bagtype string `json:"BAGTYPE"`
	// 可申购状态
	Sgzt        string `json:"SGZT"`
	Buy         bool   `json:"BUY"`
	Listtexch   string `json:"LISTTEXCH"`
	Islisttrade string `json:"ISLISTTRADE"`
}

// FundNetList 基金列表
type FundNetList []FundNetListItem

// RespFundNetList 净值列表接口原始返回结构
type RespFundNetList struct {
	Datas        []FundNetList `json:"Datas"`
	ErrCode      int           `json:"ErrCode"`
	Success      bool          `json:"Success"`
	ErrMsg       interface{}   `json:"ErrMsg"`
	Message      interface{}   `json:"Message"`
	ErrorCode    string        `json:"ErrorCode"`
	ErrorMessage interface{}   `json:"ErrorMessage"`
	ErrorMsgLst  interface{}   `json:"ErrorMsgLst"`
	TotalCount   int           `json:"TotalCount"`
	Expansion    []string      `json:"Expansion"`
}

// QueryAllFundNetList 获取天天基金净值列表全量数据，耗时很长
func (e EastMoney) QueryAllFundNetList(ctx context.Context, fundType FundType) (FundNetList, error) {
	resp, err := e.QueryFundNetListByPage(ctx, fundType, 1)
	if err != nil {
		return nil, err
	}

	result := FundNetList{}
	if len(resp.Datas) == 0 {
		return result, nil
	}
	totalCount := resp.TotalCount
	// 首页数据
	result = append(result, resp.Datas[0]...)

	// 算出总页数循环获取全量数据
	pageCount := (totalCount + 30 - 1) / 30
	for pageIndex := 2; pageIndex <= pageCount; pageIndex++ {
		resp, err := e.QueryFundNetListByPage(ctx, fundType, pageIndex)
		if err != nil {
			return result, err
		}
		if len(resp.Datas) != 0 {
			result = append(result, resp.Datas[0]...)
		}
	}
	if len(result) != totalCount {
		logging.Errorf(ctx, "QueryAllFundNetList result count:%d != resp.TotalCount:%d", len(result), totalCount)
	}
	return result, nil
}

// QueryFundNetListByPage 按页获取天天基金净值列表，接口限制最大30个
func (e EastMoney) QueryFundNetListByPage(ctx context.Context, fundType FundType, pageIndex int) (RespFundNetList, error) {
	apiurl := "https://fundmobapi.eastmoney.com/FundMNewApi/FundMNNetNewList"
	params := map[string]string{
		"FundType":   fmt.Sprint(fundType), // 基金类型
		"SortColumn": "DWJZ",               // 按净值排序
		"Sort":       "desc",               // 降序排列
		"pageIndex":  fmt.Sprint(pageIndex),
		"pageSize":   "30",
		"plat":       "Iphone",
		"deviceid":   "-",
		"product":    "EFund",
		"version":    "-",
	}
	logging.Debug(ctx, "EastMoney QueryFundNetListByPage "+apiurl+" begin", zap.Any("params", params))
	beginTime := time.Now()
	apiurl, err := goutils.NewHTTPGetURLWithQueryString(ctx, apiurl, params)
	if err != nil {
		return RespFundNetList{}, err
	}
	resp := RespFundNetList{}
	err = goutils.HTTPGET(ctx, e.HTTPClient, apiurl, &resp)
	latency := time.Now().Sub(beginTime).Milliseconds()
	logging.Debug(
		ctx,
		"EastMoney QueryFundNetListByPage "+apiurl+" end",
		zap.Int64("latency(ms)", latency),
		zap.Any("resp", resp),
	)
	if err != nil {
		return resp, err
	}
	if resp.ErrCode != 0 {
		return resp, fmt.Errorf("QueryFundNetList error %v", resp.ErrMsg)
	}
	return resp, nil
}
