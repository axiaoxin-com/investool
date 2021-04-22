// Package parser 对给定股票进行分析，得出可以买卖的股票
package parser

import (
	"context"

	"github.com/axiaoxin-com/x-stock/model"
)

// FilterStocks 筛选股票
// 1. ROE 至少 5 年内逐年递增
// 2. ROE 高于 10%
// 3. EPS 至少 5 年内逐年递增
// 4. 营业总收入至少 5 年内逐年递增
// 5. 净利润至少 5 年内逐年递增
// 6. 估值偏低或中等
// 7. 股价低于合理价格
// 8. 行业要分散
// 9. 历史波动率在 1 以内，持仓占比： 0.1:0.1-0.5:0.5-1 = 3:3:4
func FilterStocks(ctx context.Context) (model.StockList, error) {

	return nil, nil
}
