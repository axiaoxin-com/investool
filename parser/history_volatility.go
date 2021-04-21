package parser

// 历史波动率计算方法：
// 1、从市场上获得标的股票在固定时间间隔(如每天、每周或每月等)上的价格。
// 2、对于每个时间段，求出该时间段末的股价与该时段初的股价之比的自然对数。
// 3、求出这些对数值的标准差，再乘以一年中包含的时段数量的平方根(如，选取时间间隔为每天，则若扣除闭市，每年中有250个交易日，应乘以根号250)，得到的即为历史波动率。

import "context"

// HistoryVolatility 获取历史波动率
func HistoryVolatility(ctx context.Context, securityCode string) (float64, error) {
	return 0, nil
}

// RespHistoryStockPrice 历史股价接口返回结构
type RespHistoryStockPrice struct {
	Date  []string  `json:"date"`
	Price []float64 `json:"price"`
}

// QueryHistoryStockPrice 查询历史股价
func QueryHistoryStockPrice(ctx context.Context, securityCode string) (RespHistoryStockPrice, error) {
	return
}
