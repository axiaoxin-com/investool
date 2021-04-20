// 获取盈利预测

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

// RespProfitPredict 盈利预测接口返回结构
type RespProfitPredict struct {
	Version string `json:"version"`
	Result  struct {
		Pages int             `json:"pages"`
		Data  []ProfitPredict `json:"data"`
		Count int             `json:"count"`
	} `json:"result"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// ProfitPredict 盈利预测
type ProfitPredict struct {
	// 年份
	PredictYear int `json:"PREDICT_YEAR"`
	// 预测每股收益
	Eps float64 `json:"EPS"`
	// 预测市盈率
	Pe float64 `json:"PE"`
}

// QueryProfitPredict 获取盈利预测
func (e EastMoney) QueryProfitPredict(ctx context.Context, secuCode string) ([]ProfitPredict, error) {
	apiurl := "https://datacenter.eastmoney.com/securities/api/data/get"
	params := map[string]string{
		"source": "SECURITIES",
		"client": "APP",
		"type":   "RPT_RES_PROFITPREDICT",
		"sty":    "PREDICT_YEAR,EPS,PE",
		"filter": fmt.Sprintf(`(SECUCODE="%s")`, strings.ToUpper(secuCode)),
		"sr":     "1",
		"st":     "PREDICT_YEAR",
	}
	logging.Debug(ctx, "EastMoney QueryProfitPredict "+apiurl+" begin", zap.Any("params", params))
	beginTime := time.Now()
	apiurl, err := goutils.NewHTTPGetURLWithQueryString(ctx, apiurl, params)
	if err != nil {
		return nil, err
	}
	resp := RespProfitPredict{}
	if err := goutils.HTTPGET(ctx, e.HTTPClient, apiurl, &resp); err != nil {
		return nil, err
	}
	latency := time.Now().Sub(beginTime).Milliseconds()
	logging.Debug(ctx, "EastMoney QueryProfitPredict "+apiurl+" end", zap.Int64("latency(ms)", latency), zap.Any("resp", resp))
	if resp.Code != 0 {
		return nil, fmt.Errorf("%#v", resp)
	}
	return resp.Result.Data, nil
}
