// 获取估值状态

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

var (
	// ValuationLow 估值较低
	ValuationLow float64 = 0
	// ValuationModerate 估值中等
	ValuationModerate float64 = 1
	// ValuationHigh 估值较高
	ValuationHigh float64 = 2

	// ValuationMap 估值描述与数字的映射
	ValuationMap = map[string]float64{
		"估值较低": ValuationLow,
		"估值中等": ValuationModerate,
		"估值较高": ValuationHigh,
	}
)

// RespValuation 估值状态接口返回结构
type RespValuation struct {
	Version string `json:"version"`
	Result  struct {
		Pages int `json:"pages"`
		Data  []struct {
			ValationStatus string `json:"VALATION_STATUS"`
		} `json:"data"`
		Count int `json:"count"`
	} `json:"result"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// QueryValuationStatus 获取估值状态
func (e EastMoney) QueryValuationStatus(ctx context.Context, secuCode string) (float64, map[string]string, error) {
	valuations := map[string]string{}
	secuCode = strings.ToUpper(secuCode)
	apiurl := "https://datacenter.eastmoney.com/securities/api/data/get"
	// 市盈率估值
	params := map[string]string{
		"type":   "RPT_VALUATIONSTATUS",
		"sty":    "VALATION_STATUS",
		"p":      "1",
		"ps":     "1",
		"var":    "source=DataCenter",
		"client": "APP",
		"filter": fmt.Sprintf(`(SECUCODE="%s")(INDICATOR_TYPE="1")`, secuCode),
	}
	beginTime := time.Now()
	apiurl1, err := goutils.NewHTTPGetURLWithQueryString(ctx, apiurl, params)
	logging.Debug(ctx, "EastMoney QueryValuationStatus "+apiurl1+" begin", zap.Any("params", params))
	if err != nil {
		return 0, nil, err
	}
	resp := RespValuation{}
	if err := goutils.HTTPGET(ctx, e.HTTPClient, apiurl1, &resp); err != nil {
		return 0, nil, err
	}
	latency := time.Now().Sub(beginTime).Milliseconds()
	logging.Debug(
		ctx,
		"EastMoney QueryValuationStatus "+apiurl1+" end",
		zap.Int64("latency(ms)", latency),
		zap.Any("resp", resp),
	)
	if resp.Code != 0 {
		return 0, nil, fmt.Errorf("%s %#v", secuCode, resp)
	}
	if len(resp.Result.Data) > 0 {
		valuations["市盈率"] = resp.Result.Data[0].ValationStatus
	}

	// 市净率估值
	params["filter"] = fmt.Sprintf(`(SECUCODE="%s")(INDICATOR_TYPE="2")`, secuCode)
	beginTime = time.Now()
	apiurl2, err := goutils.NewHTTPGetURLWithQueryString(ctx, apiurl, params)
	logging.Debug(ctx, "EastMoney QueryValuationStatus "+apiurl2+" begin", zap.Any("params", params))
	if err != nil {
		return 0, nil, err
	}
	if err := goutils.HTTPGET(ctx, e.HTTPClient, apiurl2, &resp); err != nil {
		return 0, nil, err
	}
	latency = time.Now().Sub(beginTime).Milliseconds()
	logging.Debug(
		ctx,
		"EastMoney QueryValuationStatus "+apiurl2+" end",
		zap.Int64("latency(ms)", latency),
		zap.Any("resp", resp),
	)
	if resp.Code != 0 {
		return 0, nil, fmt.Errorf("%s %#v", secuCode, resp)
	}
	if len(resp.Result.Data) > 0 {
		valuations["市净率"] = resp.Result.Data[0].ValationStatus
	}

	// 市销率估值
	params["filter"] = fmt.Sprintf(`(SECUCODE="%s")(INDICATOR_TYPE="3")`, secuCode)
	beginTime = time.Now()
	apiurl3, err := goutils.NewHTTPGetURLWithQueryString(ctx, apiurl, params)
	logging.Debug(ctx, "EastMoney QueryValuationStatus "+apiurl3+" begin", zap.Any("params", params))
	if err != nil {
		return 0, nil, err
	}
	if err := goutils.HTTPGET(ctx, e.HTTPClient, apiurl3, &resp); err != nil {
		return 0, nil, err
	}
	latency = time.Now().Sub(beginTime).Milliseconds()
	logging.Debug(
		ctx,
		"EastMoney QueryValuationStatus "+apiurl3+" end",
		zap.Int64("latency(ms)", latency),
		zap.Any("resp", resp),
	)
	if resp.Code != 0 {
		return 0, nil, fmt.Errorf("%s %#v", secuCode, resp)
	}
	if len(resp.Result.Data) > 0 {
		valuations["市销率"] = resp.Result.Data[0].ValationStatus
	}

	// 市现率估值
	params["filter"] = fmt.Sprintf(`(SECUCODE="%s")(INDICATOR_TYPE="4")`, secuCode)
	beginTime = time.Now()
	apiurl4, err := goutils.NewHTTPGetURLWithQueryString(ctx, apiurl, params)
	logging.Debug(ctx, "EastMoney QueryValuationStatus "+apiurl4+" begin", zap.Any("params", params))
	if err != nil {
		return 0, nil, err
	}
	if err := goutils.HTTPGET(ctx, e.HTTPClient, apiurl4, &resp); err != nil {
		return 0, nil, err
	}
	latency = time.Now().Sub(beginTime).Milliseconds()
	logging.Debug(
		ctx,
		"EastMoney QueryValuationStatus "+apiurl4+" end",
		zap.Int64("latency(ms)", latency),
		zap.Any("resp", resp),
	)
	if resp.Code != 0 {
		return 0, nil, fmt.Errorf("%s %#v", secuCode, resp)
	}
	if len(resp.Result.Data) > 0 {
		valuations["市现率"] = resp.Result.Data[0].ValationStatus
	}

	// 4 种估值，综合评估一个结果
	val := float64(0)
	for _, v := range valuations {
		val += ValuationMap[v]
	}
	status := ValuationHigh
	switch {
	case val >= 0 && val < 4:
		status = ValuationLow
	case val >= 4 && val <= 6:
		status = ValuationModerate
	case val > 6:
		status = ValuationHigh
	}

	return status, valuations, nil
}
