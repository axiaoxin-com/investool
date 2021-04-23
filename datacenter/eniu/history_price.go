// 获取历史股价

package eniu

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"go.uber.org/zap"
)

// RespHistoryStockPrice 历史股价接口返回结构
type RespHistoryStockPrice struct {
	Date  []string  `json:"date"`
	Price []float64 `json:"price"`
}

// QueryHistoryStockPrice 获取历史股价
func (e Eniu) QueryHistoryStockPrice(ctx context.Context, secuCode string) (RespHistoryStockPrice, error) {
	apiurl := fmt.Sprintf("https://eniu.com/chart/pricea/%s/t/all", e.GetPathCode(ctx, secuCode))
	logging.Debug(ctx, "EastMoney QueryOrgRating "+apiurl+" begin")
	beginTime := time.Now()
	resp := RespHistoryStockPrice{}
	err := goutils.HTTPGET(ctx, e.HTTPClient, apiurl, &resp)
	latency := time.Now().Sub(beginTime).Milliseconds()
	logging.Debug(ctx, "EastMoney QueryOrgRating "+apiurl+" end", zap.Int64("latency(ms)", latency), zap.Any("resp", resp))
	return resp, err
}

// GetPathCode 返回接口 url path 中的股票代码
func (e Eniu) GetPathCode(ctx context.Context, secuCode string) string {
	s := strings.Split(secuCode, ".")
	if len(s) != 2 {
		return ""
	}
	return strings.ToLower(s[1]) + s[0]
}
