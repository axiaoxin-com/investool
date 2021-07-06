// 基金相关逻辑

package core

import (
	"context"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/x-stock/models"
	"github.com/axiaoxin-com/x-stock/services"
)

// Funder 基金实例
type Funder struct{}

// NewFunder 创建实例
func NewFunder() Funder {
	return Funder{}
}

// ParamFunderFilter Filter 参数
type ParamFunderFilter struct {
	// 类型
	Types []string `json:"types"                    form:"types"`
	// 基金规模最小值（亿）
	MinScale float64 `json:"min_scale"                form:"min_scale"`
	// 基金规模最大值（亿）
	MaxScale float64 `json:"max_scale"                form:"max_scale"`
	// 基金经理管理该基金最低年限
	MinManagerYears float64 `json:"min_manager_years"        form:"min_manager_years"`
	// 最近一年收益率排名比
	Year1RankRatio float64 `json:"year_1_rank_ratio"        form:"year_1_rank_ratio"`
	// 今年来、最近两年、最近三年、最近五年收益率排名比
	ThisYear235RankRatio float64 `json:"this_year_235_rank_ratio" form:"this_year_235_rank_ratio"`
	// 最近六月收益率排名比
	Month6RankRatio float64 `json:"month_6_rank_ratio"       form:"month_6_rank_ratio"`
	// 最近三月收益率排名比
	Month3RankRatio float64 `json:"month_3_rank_ratio"       form:"month_3_rank_ratio"`
	// 1,3,5年波动率平均值的最大值
	Max135AvgStddev float64 `json:"max_135_avg_stddev"       form:"max_135_avg_stddev"`
	// 1,3,5年夏普比率平均值的最小值
	Min135AvgSharp float64 `json:"min_135_avg_sharp"        form:"min_135_avg_sharp"`
	// 1,3,5年最大回撤率平均值的最大值
	Max135AvgRetr float64 `json:"max_135_avg_retr"         form:"max_135_avg_retr"`
	// 排序方式
	SortType models.FundSortType `json:"sort_type"                form:"sort_type"`
}

// Filter 按指定条件过滤
func (f Funder) Filter(ctx context.Context, p ParamFunderFilter) models.FundList {
	results := models.FundList{}
	for _, fund := range services.FundAllList {
		switch {
		case fund.Performance.Year5RankNum == 0:
			// 排除成立没有5年以上的基金
			continue
		case fund.Performance.Year1RankRatio > p.Year1RankRatio:
			// 最近1年排名大于前百分比的排除
			continue
		case (fund.Performance.Year2RankRatio > p.ThisYear235RankRatio || fund.Performance.Year3RankRatio > p.ThisYear235RankRatio || fund.Performance.Year5RankRatio > p.ThisYear235RankRatio || fund.Performance.ThisYearRankRatio > p.ThisYear235RankRatio):
			// 最近2,3,5以及今年来排名大于前百分比的排除
			continue
		case fund.Performance.Month6RankRatio > p.Month6RankRatio:
			// 最近6个月排名大于前百分比的排除
			continue
		case fund.Performance.Month3RankRatio > p.Month3RankRatio:
			// 最近3个月排名大于前百分比的排除
			continue
		case len(p.Types) > 0 && !goutils.IsStrInSlice(fund.Type, p.Types):
			// 指定类型时，基金类型不在指定列表则跳过
			continue
		case p.MinScale > 0 && fund.NetAssetsScale < p.MinScale*100000000:
			// 指定最小规模时，基金规模不能小于该值
			continue
		case p.MaxScale > 0 && fund.NetAssetsScale > p.MaxScale*100000000:
			// 指定最大规模时，基金规模不能大于该值
			continue
		case p.MinManagerYears > 0 && (fund.Manager.ManageDays/365) < p.MinManagerYears:
			// 指定基金经理管理该基金最低年限时，基金经理任职年数不能小于该值
			continue
		case p.Max135AvgStddev > 0:
			// 波动率平均值大于指定值时跳过
			avg, _ := goutils.AvgFloat64([]float64{fund.Stddev.Year1, fund.Stddev.Year3, fund.Stddev.Year5})
			if avg > p.Max135AvgStddev {
				continue
			}
		case p.Max135AvgRetr > 0:
			// 最大回撤率平均值大于指定值时跳过
			avg, _ := goutils.AvgFloat64(
				[]float64{fund.MaxRetracement.Year1, fund.MaxRetracement.Year3, fund.MaxRetracement.Year5},
			)
			if avg > p.Max135AvgRetr {
				continue
			}
		case p.Min135AvgSharp > 0:
			// 夏普比率平均值小于指定值时跳过
			avg, _ := goutils.AvgFloat64([]float64{fund.Sharp.Year1, fund.Sharp.Year3, fund.Sharp.Year5})
			if avg < p.Min135AvgSharp {
				continue
			}
		}
		results = append(results, fund)
	}
	results.Sort(p.SortType)
	return results
}
