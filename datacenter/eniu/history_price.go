// 获取历史股价

package eniu

import (
	"context"
	"fmt"
	"math"
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

// HistoryVolatility 按月计算历史波动率
// 历史波动率计算方法：
// 1、从市场上获得标的股票在固定时间间隔(如每天、每周或每月等)上的价格。
// 2、对于每个时间段，求出该时间段末的股价与该时段初的股价之比的自然对数。
// 3、求出这些对数值的标准差，再乘以一年中包含的时段数量的平方根(如，选取时间间隔为每天，则若扣除闭市，每年中有250个交易日，应乘以根号250)，得到的即为历史波动率。
func (p RespHistoryStockPrice) HistoryVolatility(ctx context.Context) (float64, error) {
	// 股价按月分组
	monthPrices := map[string][]float64{}
	for i, date := range p.Date {
		dayIdx := strings.LastIndex(date, "-")
		key := date[:dayIdx]
		monthPrices[key] = append(monthPrices[key], p.Price[i])
	}
	// 求末初股价比自然对数
	logs := []float64{}
	for _, prices := range monthPrices {
		startPrice := prices[0]
		endPrice := prices[len(prices)-1]
		log := math.Log(endPrice / startPrice)
		logs = append(logs, log)
	}
	// 标准差
	stddev, err := goutils.StdDeviationFloat64(logs)
	if err != nil {
		return 0, err
	}
	logging.Debugs(ctx, "stddev:", stddev)
	volatility := stddev * math.Sqrt(12)
	return volatility, nil
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
