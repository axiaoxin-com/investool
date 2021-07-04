// 基金 model

package models

import (
	"context"
	"math"
	"sort"
	"strconv"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/x-stock/datacenter/eastmoney"
)

// Fund 基金
type Fund struct {
	// 基金代码
	Code string `json:"code"`
	// 基金名称
	Name string `json:"name"`
	// 基金类型
	Type string `json:"type"`
	// 成立时间
	EstablishedDate string `json:"established_date"`
	// 最新基金净资产规模（元）
	NetAssetsScale float64 `json:"net_assets_scale"`
	// 跟踪标的代码
	IndexCode string `json:"index_code"`
	// 跟踪标的名称
	IndexName string `json:"index_name"`
	// 购买费率
	Rate string `json:"rate"`
	// 定投状态
	FixedInvestmentStatus string `json:"fixed_investment_status"`
	// 波动率
	Stddev FundStddev `json:"stddev"`
	// 最大回撤率
	MaxRetracement FundMaxRetracement `json:"max_retracement"`
	// 夏普比率
	Sharp FundSharp `json:"sharp"`
	// 绩效
	Performance FundPerformance `json:"performance"`
	// 持仓股票
	Stocks []FundStock `json:"stocks"`
	// 基金经理
	Manager FundManager `json:"manager"`
	// 历史分红送配
	HistoricalDividends []FundDividend `json:"historical_dividends"`
	// 资产占比
	AssetsProportion FundAssetsProportion `json:"assets_proportion"`
	// 行业占比
	IndustryProportions []FundIndustryProportion `json:"industry_proportions"`
}

// FundIndustryProportion 行业占比
type FundIndustryProportion struct {
	// 公布日期
	PubDate string `json:"pub_date"`
	// 行业名
	Industry string `json:"industry"`
	// 对应占比列表（%）
	Prop string `json:"prop"`
}

// FundAssetsProportion 资产占比
type FundAssetsProportion struct {
	// 公布日期
	PubDate string `json:"pub_date"`
	// 股票占比（%）
	Stock string `json:"stock"`
	// 债券占比（%）
	Bond string `json:"bond"`
	// 现金占比（%）
	Cash string `json:"cash"`
	// 其他占比（%）
	Other string `json:"other"`
	// 净资产（亿元）
	NetAssets string `json:"net_assets"`
}

// FundPerformance 基金绩效
type FundPerformance struct {
	// 近一周收益率
	WeekProfitRatio float64 `json:"week_profit_ratio"`
	// 近一周涨跌幅
	WeekAmplitude float64 `json:"week_amplitude"`
	// 近一周同类均值
	WeekKindAvg float64 `json:"week_kind_avg"`
	// 近一周同类排名
	WeekRankNum float64 `json:"week_rank_num"`
	// 近一周同类排名百分比
	WeekRankRatio float64 `json:"week_rank_ratio"`
	// 近一周同类总数
	WeekRankTotalCount float64 `json:"week_rank_total_count"`
	// 近一月收益率
	Month1ProfitRatio float64 `json:"month_1_profit_ratio"`
	// 近一月涨跌幅
	Month1Amplitude float64 `json:"month_1_amplitude"`
	// 近一月同类均值
	Month1KindAvg float64 `json:"month_1_kind_avg"`
	// 近一月同类排名
	Month1RankNum float64 `json:"month_1_rank_num"`
	// 近一月同类排名百分比
	Month1RankRatio float64 `json:"month_1_rank_ratio"`
	// 近一月同类总数
	Month1RankTotalCount float64 `json:"month_1_rank_total_count"`
	// 近三月收益率
	Month3ProfitRatio float64 `json:"month_3_profit_ratio"`
	// 近三月涨跌幅
	Month3Amplitude float64 `json:"month_3_amplitude"`
	// 近三月同类均值
	Month3KindAvg float64 `json:"month_3_kind_avg"`
	// 近三月同类排名
	Month3RankNum float64 `json:"month_3_rank_num"`
	// 近三月同类排名百分比
	Month3RankRatio float64 `json:"month_3_rank_ratio"`
	// 近三月同类总数
	Month3RankTotalCount float64 `json:"month_3_rank_total"`
	// 近六月收益率
	Month6ProfitRatio float64 `json:"month_6_profit_ratio"`
	// 近六月涨跌幅
	Month6Amplitude float64 `json:"month_6_amplitude"`
	// 近六月同类均值
	Month6KindAvg float64 `json:"month_6_kind_avg"`
	// 近六月同类排名
	Month6RankNum float64 `json:"month_6_rank_num"`
	// 近六月同类排名百分比
	Month6RankRatio float64 `json:"month_6_rank_ratio"`
	// 近六月同类总数
	Month6RankTotalCount float64 `json:"month_6_rank_total_count"`
	// 近一年收益率
	Year1ProfitRatio float64 `json:"year_1_profit_ratio"`
	// 近一年涨跌幅
	Year1Amplitude float64 `json:"year_1_amplitude"`
	// 近一年同类均值
	Year1KindAvg float64 `json:"year_1_kind_avg"`
	// 近一年同类排名
	Year1RankNum float64 `json:"year_1_rank_num"`
	// 近一年同类排名百分比
	Year1RankRatio float64 `json:"year_1_rank_ratio"`
	// 近一年同类总数
	Year1RankTotalCount float64 `json:"year_1_rank_total_count"`
	// 近两年收益率
	Year2ProfitRatio float64 `json:"year_2_profit_ratio"`
	// 近两年涨跌幅
	Year2Amplitude float64 `json:"year_2_amplitude"`
	// 近两年同类均值
	Year2KindAvg float64 `json:"year_2_kind_avg"`
	// 近两年同类排名
	Year2RankNum float64 `json:"year_2_rank_num"`
	// 近两年同类排名百分比
	Year2RankRatio float64 `json:"year_2_rank_ratio"`
	// 近两年同类总数
	Year2RankTotalCount float64 `json:"year_2_rank_total_count"`
	// 近三年收益率
	Year3ProfitRatio float64 `json:"year_3_profit_ratio"`
	// 近三年涨跌幅
	Year3Amplitude float64 `json:"year_3_amplitude"`
	// 近三年同类均值
	Year3KindAvg float64 `json:"year_3_kind_avg"`
	// 近三年同类排名
	Year3RankNum float64 `json:"year_3_rank_num"`
	// 近三年同类排名百分比
	Year3RankRatio float64 `json:"year_3_rank_ratio"`
	// 近三年同类总数
	Year3RankTotalCount float64 `json:"year_3_rank_total_count"`
	// 近五年收益率
	Year5ProfitRatio float64 `json:"year_5_profit_ratio"`
	// 近五年涨跌幅
	Year5Amplitude float64 `json:"year_5_amplitude"`
	// 近五年同类均值
	Year5KindAvg float64 `json:"year_5_kind_avg"`
	// 近五年同类排名
	Year5RankNum float64 `json:"year_5_rank_num"`
	// 近五年同类排名百分比
	Year5RankRatio float64 `json:"year_5_rank_ratio"`
	// 近五年同类总数
	Year5RankTotalCount float64 `json:"year_5_rank_total_count"`
	// 今年来收益率
	ThisYearProfitRatio float64 `json:"this_year_profit_ratio"`
	// 今年来涨跌幅
	ThisYearAmplitude float64 `json:"this_year_amplitude"`
	// 今年来同类均值
	ThisYearKindAvg float64 `json:"this_year_kind_avg"`
	// 今年来同类排名
	ThisYearRankNum float64 `json:"this_year_rank_num"`
	// 今年来同类排名百分比
	ThisYearRankRatio float64 `json:"this_year_rank_ratio"`
	// 今年来同类总数
	ThisYearRankTotalCount float64 `json:"this_year_rank_total_count"`
	// 成立以来收益率
	HistoricalProfitRatio float64 `json:"historical_profit_ratio"`
	// 成立以来涨跌幅
	HistoricalAmplitude float64 `json:"historical_amplitude"`
	// 成立以来同类均值
	HistoricalKindAvg float64 `json:"historical_kind_avg"`
	// 成立以来同类排名
	HistoricalRankNum float64 `json:"historical_rank_num"`
	// 成立以来同类排名百分比
	HistoricalRankRatio float64 `json:"historical_rank_ratio"`
	// 成立以来同类总数
	HistoricalRankTotalCount float64 `json:"historical_rank_total_count"`
}

// FundDividend 分红送配
type FundDividend struct {
	// 权益登记日
	RegDate string `json:"reg_date"`
	// 每份分红（元）
	Value float64 `json:"value"`
	// 分红发放日
	RationDate string `json:"ration_date"`
}

// FundStddev 波动率
type FundStddev struct {
	// 近1年波动率（%）
	Year1 float64 `json:"year_1"`
	// 近3年波动率（%）
	Year3 float64 `json:"year_3"`
	// 近5年波动率（%）
	Year5 float64 `json:"year_5"`
}

// FundMaxRetracement 最大回撤
type FundMaxRetracement struct {
	// 近1年最大回撤（%）
	Year1 float64 `json:"year_1"`
	// 近3年最大回撤（%）
	Year3 float64 `json:"year_3"`
	// 近5年最大回撤（%）
	Year5 float64 `json:"year_5"`
}

// FundSharp 夏普比率
type FundSharp struct {
	// 近1年夏普比率
	Year1 float64 `json:"year_1"`
	// 近3年夏普比率
	Year3 float64 `json:"year_3"`
	// 近5年夏普比率
	Year5 float64 `json:"year_5"`
}

// FundStock 基金持仓股票
type FundStock struct {
	// 股票代码
	Code string `json:"code"`
	// 股票名称
	Name string `json:"name"`
	// 交易所代号
	ExCode string `json:"ex_code"`
	// 股票行业
	Industry string `json:"industry"`
	// 持仓占比(%)
	HoldRatio float64 `json:"hold_ratio"`
	// 较上期调仓比例
	AdjustRatio float64 `json:"adjust_ratio"`
}

// FundManager 基金经理
type FundManager struct {
	// ID
	ID string `json:"id"`
	// 基金经理名字
	Name string `json:"name"`
	// 从业时间（天）
	WorkingDays float64 `json:"working_days"`
	// 管理该基金时间（天）
	ManageDays float64 `json:"manage_days"`
	// 任职回报（%）
	ManageRepay float64 `json:"manage_repay"`
	// 年均回报（%）
	YearsAvgRepay float64 `json:"years_avg_repay"`
}

func interfaceToFloat64(ctx context.Context, unk interface{}) (result float64) {
	var err error
	switch i := unk.(type) {
	case float64:
		result = float64(i)
	case float32:
		result = float64(i)
	case int64:
		result = float64(i)
	case int32:
		result = float64(i)
	case int:
		result = float64(i)
	case uint32:
		result = float64(i)
	case uint64:
		result = float64(i)
	case uint:
		result = float64(i)
	case string:
		if i == "" || i == "--" {
			result = 0.0
		} else {
			result, err = strconv.ParseFloat(i, 64)
			if err != nil {
				logging.Errorf(ctx, "interfaceToFloat64 ParseFloat error:%v i:%v unk:%v", err, i, unk)
			}
		}
	default:
		result = 0.0
	}
	return
}

// NewFund 创建 Fund 实例
func NewFund(ctx context.Context, efund eastmoney.RespFundInfo) Fund {
	fund := Fund{
		Code:            efund.Jjxq.Datas.Fcode,
		Name:            efund.Jjxq.Datas.Shortname,
		Type:            efund.Jjxq.Datas.Ftype,
		EstablishedDate: efund.Jjxq.Datas.Estabdate,
		IndexCode:       efund.Jjxq.Datas.Indexcode,
		IndexName:       efund.Jjxq.Datas.Indexname,
		Rate:            efund.Jjxq.Datas.Rate,
		Stddev: FundStddev{
			Year1: interfaceToFloat64(ctx, efund.Tssj.Datas.Stddev1),
			Year3: interfaceToFloat64(ctx, efund.Tssj.Datas.Stddev3),
			Year5: interfaceToFloat64(ctx, efund.Tssj.Datas.Stddev5),
		},

		MaxRetracement: FundMaxRetracement{
			Year1: interfaceToFloat64(ctx, efund.Tssj.Datas.Maxretra1),
			Year3: interfaceToFloat64(ctx, efund.Tssj.Datas.Maxretra3),
			Year5: interfaceToFloat64(ctx, efund.Tssj.Datas.Maxretra5),
		},
		Sharp: FundSharp{
			Year1: interfaceToFloat64(ctx, efund.Tssj.Datas.Sharp1),
			Year3: interfaceToFloat64(ctx, efund.Tssj.Datas.Sharp3),
			Year5: interfaceToFloat64(ctx, efund.Tssj.Datas.Sharp5),
		},
	}

	// 定投状态
	switch efund.Jjxq.Datas.Dtzt {
	case "1":
		fund.FixedInvestmentStatus = "可定投"
	}

	// 基金规模
	if len(efund.Jjgm.Datas) > 0 {
		fund.NetAssetsScale = interfaceToFloat64(ctx, efund.Jjgm.Datas[0].Netnav)
	} else {
		logging.Warnf(ctx, "code:%v jjgm no data", fund.Code)
	}

	// 绩效
	pfm := FundPerformance{}
	for _, d := range efund.Jdzf.Datas {
		rankRatio := interfaceToFloat64(ctx, interfaceToFloat64(ctx, d.Rank)) / interfaceToFloat64(ctx, d.Sc)
		if math.IsNaN(rankRatio) {
			rankRatio = 0.0
		}
		rankRatio = rankRatio * 100.0 // %
		switch d.Title {
		case "Z":
			pfm.WeekAmplitude = interfaceToFloat64(ctx, d.Avg)
			pfm.WeekKindAvg = interfaceToFloat64(ctx, d.Hs300)
			pfm.WeekRankNum = interfaceToFloat64(ctx, d.Rank)
			pfm.WeekRankRatio = rankRatio
			pfm.WeekProfitRatio = interfaceToFloat64(ctx, d.Syl)
			pfm.WeekRankTotalCount = interfaceToFloat64(ctx, d.Sc)
		case "Y":
			pfm.Month1Amplitude = interfaceToFloat64(ctx, d.Avg)
			pfm.Month1KindAvg = interfaceToFloat64(ctx, d.Hs300)
			pfm.Month1RankNum = interfaceToFloat64(ctx, d.Rank)
			pfm.Month1RankRatio = rankRatio
			pfm.Month1ProfitRatio = interfaceToFloat64(ctx, d.Syl)
			pfm.Month1RankTotalCount = interfaceToFloat64(ctx, d.Sc)
		case "3Y":
			pfm.Month3Amplitude = interfaceToFloat64(ctx, d.Avg)
			pfm.Month3KindAvg = interfaceToFloat64(ctx, d.Hs300)
			pfm.Month3RankNum = interfaceToFloat64(ctx, d.Rank)
			pfm.Month3RankRatio = rankRatio
			pfm.Month3ProfitRatio = interfaceToFloat64(ctx, d.Syl)
			pfm.Month3RankTotalCount = interfaceToFloat64(ctx, d.Sc)
		case "6Y":
			pfm.Month6Amplitude = interfaceToFloat64(ctx, d.Avg)
			pfm.Month6KindAvg = interfaceToFloat64(ctx, d.Hs300)
			pfm.Month6RankNum = interfaceToFloat64(ctx, d.Rank)
			pfm.Month6RankRatio = rankRatio
			pfm.Month6ProfitRatio = interfaceToFloat64(ctx, d.Syl)
			pfm.Month6RankTotalCount = interfaceToFloat64(ctx, d.Sc)
		case "1N":
			pfm.Year1Amplitude = interfaceToFloat64(ctx, d.Avg)
			pfm.Year1KindAvg = interfaceToFloat64(ctx, d.Hs300)
			pfm.Year1RankNum = interfaceToFloat64(ctx, d.Rank)
			pfm.Year1RankRatio = rankRatio
			pfm.Year1ProfitRatio = interfaceToFloat64(ctx, d.Syl)
			pfm.Year1RankTotalCount = interfaceToFloat64(ctx, d.Sc)
		case "2N":
			pfm.Year2Amplitude = interfaceToFloat64(ctx, d.Avg)
			pfm.Year2KindAvg = interfaceToFloat64(ctx, d.Hs300)
			pfm.Year2RankNum = interfaceToFloat64(ctx, d.Rank)
			pfm.Year2RankRatio = rankRatio
			pfm.Year2ProfitRatio = interfaceToFloat64(ctx, d.Syl)
			pfm.Year2RankTotalCount = interfaceToFloat64(ctx, d.Sc)
		case "3N":
			pfm.Year3Amplitude = interfaceToFloat64(ctx, d.Avg)
			pfm.Year3KindAvg = interfaceToFloat64(ctx, d.Hs300)
			pfm.Year3RankNum = interfaceToFloat64(ctx, d.Rank)
			pfm.Year3RankRatio = rankRatio
			pfm.Year3ProfitRatio = interfaceToFloat64(ctx, d.Syl)
			pfm.Year3RankTotalCount = interfaceToFloat64(ctx, d.Sc)
		case "5N":
			pfm.Year5Amplitude = interfaceToFloat64(ctx, d.Avg)
			pfm.Year5KindAvg = interfaceToFloat64(ctx, d.Hs300)
			pfm.Year5RankNum = interfaceToFloat64(ctx, d.Rank)
			pfm.Year5RankRatio = rankRatio
			pfm.Year5ProfitRatio = interfaceToFloat64(ctx, d.Syl)
			pfm.Year5RankTotalCount = interfaceToFloat64(ctx, d.Sc)
		case "JN":
			pfm.ThisYearAmplitude = interfaceToFloat64(ctx, d.Avg)
			pfm.ThisYearKindAvg = interfaceToFloat64(ctx, d.Hs300)
			pfm.ThisYearRankNum = interfaceToFloat64(ctx, d.Rank)
			pfm.ThisYearRankRatio = rankRatio
			pfm.ThisYearProfitRatio = interfaceToFloat64(ctx, d.Syl)
			pfm.ThisYearRankTotalCount = interfaceToFloat64(ctx, d.Sc)
		case "LN":
			pfm.HistoricalAmplitude = interfaceToFloat64(ctx, d.Avg)
			pfm.HistoricalKindAvg = interfaceToFloat64(ctx, d.Hs300)
			pfm.HistoricalRankNum = interfaceToFloat64(ctx, d.Rank)
			pfm.HistoricalRankRatio = rankRatio
			pfm.HistoricalProfitRatio = interfaceToFloat64(ctx, d.Syl)
			pfm.HistoricalRankTotalCount = interfaceToFloat64(ctx, d.Sc)
		}
	}
	fund.Performance = pfm

	// 持仓股票
	stocks := []FundStock{}
	for _, s := range efund.Jjcc.Datas.InverstPosition.FundStocks {
		stock := FundStock{
			Code:        s.Gpdm,
			Name:        s.Gpjc,
			Industry:    s.Indexname,
			ExCode:      s.Newtexch,
			HoldRatio:   interfaceToFloat64(ctx, s.Jzbl),
			AdjustRatio: interfaceToFloat64(ctx, s.Pctnvchg),
		}
		stocks = append(stocks, stock)
	}
	fund.Stocks = stocks

	// 基金经理
	manager := FundManager{}
	if len(efund.Jjjlnew.Datas) > 0 {
		jjjl := efund.Jjjlnew.Datas[0]
		if len(jjjl.Manger) > 0 {
			m := jjjl.Manger[0]
			manager.ID = m.Mgrid
			manager.Name = m.Mgrname
			manager.WorkingDays = interfaceToFloat64(ctx, m.Totaldays)
			manager.ManageDays = interfaceToFloat64(ctx, jjjl.Days)
			manager.ManageRepay = interfaceToFloat64(ctx, jjjl.Penavgrowth)
			manager.YearsAvgRepay = interfaceToFloat64(ctx, m.Yieldse)
			fund.Manager = manager
		} else {
			logging.Warnf(ctx, "code:%v jjjlnew manager no data", fund.Code)
		}
	} else {
		logging.Warnf(ctx, "code:%v jjjlnew no data", fund.Code)
	}

	// 分红送配
	dividends := []FundDividend{}
	for _, d := range efund.Fhsp.Datas.Fhinfo {
		fd := FundDividend{
			RegDate:    d.Djr,
			Value:      interfaceToFloat64(ctx, d.Fhfcz),
			RationDate: d.Ffr,
		}
		dividends = append(dividends, fd)
	}
	if len(dividends) > 5 {
		dividends = dividends[:5]
	}
	fund.HistoricalDividends = dividends

	// 资产占比
	for _, vlist := range efund.Jjcc.Datas.AssetAllocation {
		if len(vlist) > 0 {
			v := vlist[0]
			ap := FundAssetsProportion{
				PubDate:   v["FSRQ"],
				Stock:     v["GP"] + "%",
				Bond:      v["ZQ"] + "%",
				Cash:      v["HB"] + "%",
				Other:     v["QT"] + "%",
				NetAssets: v["JZC"] + "亿",
			}
			fund.AssetsProportion = ap
		}
	}

	// 行业占比
	for date, vlist := range efund.Jjcc.Datas.SectorAllocation {
		for _, i := range vlist {
			if i["ZJZBL"] == "0" || i["ZJZBL"] == "--" {
				continue
			}
			ip := FundIndustryProportion{
				PubDate:  date,
				Industry: i["HYMC"],
				Prop:     i["ZJZBL"],
			}
			fund.IndustryProportions = append(fund.IndustryProportions, ip)
		}
	}

	return fund
}

// FundList list
type FundList []Fund

// FundSortType 基金排序类型
type FundSortType int

const (
	// FundSortTypeWeek 按最近一周收益率排序
	FundSortTypeWeek = iota
	// FundSortTypeMonth1 按最近一月收益率排序
	FundSortTypeMonth1
	// FundSortTypeMonth3 按最新三月收益率排序
	FundSortTypeMonth3
	// FundSortTypeMonth6 按最新六月收益率排序
	FundSortTypeMonth6
	// FundSortTypeYear1 按最近一年收益率排序
	FundSortTypeYear1
	// FundSortTypeYear2 按最近两年收益率排序
	FundSortTypeYear2
	// FundSortTypeYear3 按最近三年收益率排序
	FundSortTypeYear3
	// FundSortTypeYear5 按最近五年收益率排序
	FundSortTypeYear5
	// FundSortTypeThisYear 按今年来收益率排序
	FundSortTypeThisYear
	// FundSortTypeHistorical 按成立来收益率排序
	FundSortTypeHistorical
)

// Sort 排序
func (f FundList) Sort(st FundSortType) {
	switch st {
	case FundSortTypeWeek:
		sort.Slice(f, func(i, j int) bool {
			return f[i].Performance.WeekProfitRatio > f[j].Performance.WeekProfitRatio
		})
	case FundSortTypeMonth1:
		sort.Slice(f, func(i, j int) bool {
			return f[i].Performance.Month1ProfitRatio > f[j].Performance.Month1ProfitRatio
		})
	case FundSortTypeMonth3:
		sort.Slice(f, func(i, j int) bool {
			return f[i].Performance.Month3ProfitRatio > f[j].Performance.Month3ProfitRatio
		})
	case FundSortTypeMonth6:
		sort.Slice(f, func(i, j int) bool {
			return f[i].Performance.Month6ProfitRatio > f[j].Performance.Month6ProfitRatio
		})
	case FundSortTypeYear1:
		sort.Slice(f, func(i, j int) bool {
			return f[i].Performance.Year1ProfitRatio > f[j].Performance.Year1ProfitRatio
		})
	case FundSortTypeYear2:
		sort.Slice(f, func(i, j int) bool {
			return f[i].Performance.Year2ProfitRatio > f[j].Performance.Year2ProfitRatio
		})
	case FundSortTypeYear3:
		sort.Slice(f, func(i, j int) bool {
			return f[i].Performance.Year3ProfitRatio > f[j].Performance.Year3ProfitRatio
		})
	case FundSortTypeYear5:
		sort.Slice(f, func(i, j int) bool {
			return f[i].Performance.Year5ProfitRatio > f[j].Performance.Year5ProfitRatio
		})
	case FundSortTypeThisYear:
		sort.Slice(f, func(i, j int) bool {
			return f[i].Performance.ThisYearProfitRatio > f[j].Performance.ThisYearProfitRatio
		})
	case FundSortTypeHistorical:
		sort.Slice(f, func(i, j int) bool {
			return f[i].Performance.HistoricalProfitRatio > f[j].Performance.HistoricalProfitRatio
		})
	}
}

// Is4433 判断是否满足4433法则
func (f Fund) Is4433(ctx context.Context) bool {
	// 没有5年数据则不满足
	if f.Performance.Year5ProfitRatio == 0 || f.Performance.Year5RankNum == 0 {
		return false
	}
	quarterRatio := float64(1) / float64(4) * 100.0
	oneThirdRatio := float64(1) / float64(3) * 100.0
	// 最近1年收益率排名在同类型基金的前四分之一；
	if f.Performance.Year1RankRatio > quarterRatio {
		return false
	}
	// 最近2年、3年、5年及今年以来收益率排名均在同类型基金的前四分之一；
	if f.Performance.Year2RankRatio > quarterRatio || f.Performance.Year3RankRatio > quarterRatio || f.Performance.Year5RankRatio > quarterRatio ||
		f.Performance.ThisYearRankRatio > quarterRatio {
		return false
	}
	// 最近6个月收益率排名在同类型基金的前三分之一；
	if f.Performance.Month6RankRatio > oneThirdRatio {
		return false
	}
	// 最近3个月收益率排名在同类型基金的前三分之一；
	if f.Performance.Month3RankRatio > oneThirdRatio {
		return false
	}
	return true
}

// NetAssetsScaleHuman 净资产数字转换为亿、万单位
func (f Fund) NetAssetsScaleHuman() string {
	return goutils.YiWanString(f.NetAssetsScale)
}
