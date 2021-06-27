// 基金相关逻辑

package core

import (
	"context"

	"github.com/axiaoxin-com/x-stock/models"
	"github.com/axiaoxin-com/x-stock/services"
)

// Funder 基金实例
type Funder struct{}

// NewFunder 创建实例
func NewFunder() Funder {
	return Funder{}
}

// ParamFilterByRankRatio FilterByRankRatio 参数
type ParamFilterByRankRatio struct {
	// 最近一周收益率排名比
	Week float64
	// 最近一月收益率排名比
	Month1 float64
	// 最近三月收益率排名比
	Month3 float64
	// 最近六月收益率排名比
	Month6 float64
	// 今年来收益率排名比
	ThisYear float64
	// 最近一年收益率排名比
	Year1 float64
	// 最近两年收益率排名比
	Year2 float64
	// 最近三年收益率排名比
	Year3 float64
	// 最近五年收益率排名比
	Year5 float64
}

// FilterByRankRatio 按指定收益率排名前百分比过滤
func (f Funder) FilterByRankRatio(ctx context.Context, p ParamFilterByRankRatio) models.FundList {
	results := models.FundList{}
	for _, fund := range services.FundAllList {
		if fund.Performance.WeekProfitRatio <= p.Week &&
			fund.Performance.Month1RankRatio <= p.Month1 &&
			fund.Performance.Month6RankRatio <= p.Month6 &&
			fund.Performance.Year1ProfitRatio <= p.Year1 &&
			fund.Performance.Year2ProfitRatio <= p.Year2 &&
			fund.Performance.Year3ProfitRatio <= p.Year3 &&
			fund.Performance.Year5ProfitRatio <= p.Year5 &&
			fund.Performance.ThisYearRankRatio <= p.ThisYear {
			results = append(results, fund)
		}
	}
	results.SortByYear1RankRatio()
	return results
}
