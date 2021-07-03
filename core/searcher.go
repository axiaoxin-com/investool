// 关键词搜索股票

package core

import (
	"context"
	"fmt"

	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/x-stock/datacenter"
	"github.com/axiaoxin-com/x-stock/datacenter/eastmoney"
	"github.com/axiaoxin-com/x-stock/datacenter/qq"
	"github.com/axiaoxin-com/x-stock/models"
)

// Searcher 搜索器实例
type Searcher struct{}

// NewSearcher 创建搜索器实例
func NewSearcher(ctx context.Context) Searcher {
	return Searcher{}
}

// Search 按股票名或代码搜索股票
func (c Searcher) Search(ctx context.Context, keywords []string) (models.StockList, error) {
	// 根据关键词匹配股票代码
	matchedResults := []qq.SearchResult{}
	for _, kw := range keywords {
		searchResults, err := datacenter.QQ.KeywordSearch(ctx, kw)
		if err != nil {
			logging.Errorf(ctx, "search %s error:", kw, err.Error())
			continue
		}
		if len(searchResults) == 0 {
			logging.Warnf(ctx, "search %s no data", kw)
			continue
		}
		logging.Debugf(ctx, "search results:%+v, %s matched", searchResults, searchResults[0])
		matchedResults = append(matchedResults, searchResults[0])
	}
	if len(matchedResults) == 0 {
		return nil, fmt.Errorf("can't find the %v", keywords)
	}
	// 查询匹配到的股票代码的股票信息
	filter := eastmoney.Filter{}
	for _, result := range matchedResults {
		filter.SpecialSecurityCodeList = append(filter.SpecialSecurityCodeList, result.SecurityCode)
	}
	stocks, err := datacenter.EastMoney.QuerySelectedStocksWithFilter(ctx, filter)
	if err != nil {
		return nil, err
	}
	results := models.StockList{}
	for _, stock := range stocks {
		mstock, err := models.NewStock(ctx, stock)
		if err != nil {
			logging.Errorf(ctx, "%s new models stock error:%v", stock.SecurityCode, err.Error())
			continue
		}
		results = append(results, mstock)
	}
	return results, nil
}
