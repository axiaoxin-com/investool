// 关键词搜索股票

package core

import (
	"context"

	"github.com/axiaoxin-com/x-stock/model"
)

// Searcher 搜索器实例
type Searcher struct{}

// NewSearcher 创建搜索器实例
func NewSearcher(ctx context.Context, stock model.Stock) Searcher {
	return Searcher{}
}

// Search 按股票名或代码搜索股票
func (c Searcher) Search(ctx context.Context, keyword string) (results model.StockList, err error) {
	return
}
