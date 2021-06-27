// 基金 model

package models

import (
	"errors"

	"github.com/axiaoxin-com/x-stock/datacenter/eastmoney"
)

// Fund 基金
type Fund struct {
	// 基金代码
	Code string
	// 基金名称
	Name string
	// 基金类型
	Type string
	// 成立时间
	EstablishedDate string
	// 最新基金净资产规模（元）
	NetAssetsScale float64
	// 跟踪标的代码
	IndexCode string
	// 跟踪标的名称
	IndexName string
	// 购买费率
	Rate float64
	// 定投状态
	FixedInvestmentStatus string
	// 波动率
	Stddev FundStddev
	// 最大回撤率
	MaxRetracement FundMaxRetracement
	Sharp          FundSharp
	// 绩效
	Performance FundPerformance
	// 持仓股票
	Stocks []FundStock
	// 基金经理
	Manager FundManager
	// 历史分红送配
	HistoricalDividends []FundDividend
	// 资产占比
	AssetsProportion FundAssetsProportion
	// 行业占比
	IndustryProportion FundIndustryProportion
}

// FundIndustryProportion 行业占比
type FundIndustryProportion struct {
	// 公布日期
	PubDate string
	// 行业名称列表
	Industry []string
	// 对应占比列表（%）
	Props []string
}

// FundAssetsProportion 资产占比
type FundAssetsProportion struct {
	// 公布日期
	PubDate string
	// 股票占比（%）
	Stock string
	// 债券占比（%）
	Bond string
	// 现金占比（%）
	Cash string
	// 其他占比（%）
	Other string
	// 净资产（亿元）
	NetAssets string
}

// FundPerformance 基金绩效
type FundPerformance struct {
	// 今年来收益率
	ThisYearProfitRatio float64
	// 今年来涨跌幅
	ThisYearAmplitude float64
	// 今年来同类均值
	ThisYearKindAvg float64
	// 今年来同类排名
	ThisYearRankNum int
	// 今年来同类排名百分比
	ThisYearRankRatio float64
	// 成立以来收益率
	HistoricalProfitRatio float64
	// 成立以来涨跌幅
	HistoricalAmplitude float64
	// 成立以来同类均值
	HistoricalKindAvg float64
	// 成立以来同类排名
	HistoricalRankNum int
	// 成立以来同类排名百分比
	HistoricalRankRatio float64
	// 近一周收益率
	WeekProfitRatio float64
	// 近一周涨跌幅
	WeekAmplitude float64
	// 近一周同类均值
	WeekKindAvg float64
	// 近一周同类排名
	WeekRankNum int
	// 近一周同类排名百分比
	WeekRankRatio float64
	// 近一月收益率
	Month1ProfitRatio float64
	// 近一月涨跌幅
	Month1Amplitude float64
	// 近一月同类均值
	Month1KindAvg float64
	// 近一月同类排名
	Month1RankNum int
	// 近一月同类排名百分比
	Month1RankRatio float64
	// 近三月收益率
	Month3ProfitRatio float64
	// 近三月涨跌幅
	Month3Amplitude float64
	// 近三月同类均值
	Month3KindAvg float64
	// 近三月同类排名
	Month3RankNum int
	// 近三月同类排名百分比
	Month3RankRatio float64
	// 近六月收益率
	Month6ProfitRatio float64
	// 近六月涨跌幅
	Month6Amplitude float64
	// 近六月同类均值
	Month6KindAvg float64
	// 近六月同类排名
	Month6RankNum int
	// 近六月同类排名百分比
	Month6RankRatio float64
	// 近一年收益率
	Year1ProfitRatio float64
	// 近一年涨跌幅
	Year1Amplitude float64
	// 近一年同类均值
	Year1KindAvg float64
	// 近一年同类排名
	Year1RankNum int
	// 近一年同类排名百分比
	Year1RankRatio float64
	// 近两年收益率
	Year2ProfitRatio float64
	// 近两年涨跌幅
	Year2Amplitude float64
	// 近两年同类均值
	Year2KindAvg float64
	// 近两年同类排名
	Year2RankNum int
	// 近两年同类排名百分比
	Year2RankRatio float64
	// 近三年收益率
	Year3ProfitRatio float64
	// 近三年涨跌幅
	Year3Amplitude float64
	// 近三年同类均值
	Year3KindAvg float64
	// 近三年同类排名
	Year3RankNum int
	// 近三年同类排名百分比
	Year3RankRatio float64
	// 近五年收益率
	Year5ProfitRatio float64
	// 近五年涨跌幅
	Year5Amplitude float64
	// 近五年同类均值
	Year5KindAvg float64
	// 近五年同类排名
	Year5RankNum int
	// 近五年同类排名百分比
	Year5RankRatio float64
}

// FundDividend 分红送配
type FundDividend struct {
	// 权益登记日
	RegDate string
	// 每份分红（元）
	Value float64
	// 分红发放日
	RationDate string
}

// FundStddev 波动率
type FundStddev struct {
	// 近1年波动率（%）
	Year1 float64
	// 近3年波动率（%）
	Year3 float64
	// 近5年波动率（%）
	Year5 float64
}

// FundMaxRetracement 最大回撤
type FundMaxRetracement struct {
	// 近1年最大回撤（%）
	Year1 float64
	// 近3年最大回撤（%）
	Year3 float64
	// 近5年最大回撤（%）
	Year5 float64
}

// FundSharp 夏普比率
type FundSharp struct {
	// 近1年夏普比率
	Year1 float64
	// 近3年夏普比率
	Year3 float64
	// 近5年夏普比率
	Year5 float64
}

// FundStock 基金持仓股票
type FundStock struct {
	// 股票代码
	Code string
	// 股票名称
	Name string
	// 股票行业
	Industry string
	// 持仓占比(%)
	HoldRatio float64
	// 较上期调仓比例
	AdjustRatio float64
}

// FundManager 基金经理
type FundManager struct {
	// 基金经理名字
	Name string
	// 从业时间（天）
	WorkingDays int
	// 管理该基金时间（天）
	ManageDays int
	// 任职回报（%）
	ManageRepay float64
	// 年均回报（%）
	YearsAvgRepay float64
}

// NewFund 创建 Fund 实例
func NewFund(efund eastmoney.RespFundInfo) (Fund, error) {
	fund := Fund{
		Code:            efund.Jjxq.Datas.Fcode,
		Name:            efund.Jjxq.Datas.Shortname,
		Type:            efund.Jjxq.Datas.Ftype,
		EstablishedDate: efund.Jjxq.Datas.Estabdate,
		IndexCode:       efund.Jjxq.Datas.Indexcode,
		IndexName:       efund.Jjxq.Datas.Indexname,
		Rate:            efund.Jjxq.Datas.Rate,
		Stddev: FundStddev{
			Year1: efund.Tssj.Datas.Stddev1,
			Year3: efund.Tssj.Datas.Stddev3,
			Year5: efund.Tssj.Datas.Stddev5,
		},

		MaxRetracement: FundMaxRetracement{
			Year1: efund.Tssj.Datas.Maxretra1,
			Year3: efund.Tssj.Datas.Maxretra3,
			Year5: efund.Tssj.Datas.Maxretra5,
		},
		Sharp: FundSharp{
			Year1: efund.Tssj.Datas.Sharp1,
			Year3: efund.Tssj.Datas.Sharp3,
			Year5: efund.Tssj.Datas.Sharp5,
		},
	}

	// 定投状态
	switch efund.Jjxq.Datas.Dtzt {
	case "1":
		fund.FixedInvestmentStatus = "可定投"
	}

	// 基金规模
	if len(efund.Jjgm.Datas) == 0 {
		return fund, errors.New("jjgm no data")
	}
	jjgm := efund.Jjgm.Datas[0]
	fund.NetAssetsScale = jjgm.Netnav

	// 绩效
	pfm := FundPerformance{}
	for _, d := range efund.Jdzf.Datas {
		rankRatio := float64(d.Rank) / float64(d.Sc)
		switch d.Title {
		case "Z":
			pfm.WeekAmplitude = d.Avg
			pfm.WeekKindAvg = d.Hs300
			pfm.WeekRankNum = d.Rank
			pfm.WeekRankRatio = rankRatio
			pfm.WeekProfitRatio = d.Syl
		case "Y":
			pfm.Month1Amplitude = d.Avg
			pfm.Month1KindAvg = d.Hs300
			pfm.Month1RankNum = d.Rank
			pfm.Month1RankRatio = rankRatio
			pfm.Month1ProfitRatio = d.Syl
		case "3Y":
			pfm.Month3Amplitude = d.Avg
			pfm.Month3KindAvg = d.Hs300
			pfm.Month3RankNum = d.Rank
			pfm.Month3RankRatio = rankRatio
			pfm.Month3ProfitRatio = d.Syl
		case "6Y":
			pfm.Month6Amplitude = d.Avg
			pfm.Month6KindAvg = d.Hs300
			pfm.Month6RankNum = d.Rank
			pfm.Month6RankRatio = rankRatio
			pfm.Month6ProfitRatio = d.Syl
		case "1N":
			pfm.Year1Amplitude = d.Avg
			pfm.Year1KindAvg = d.Hs300
			pfm.Year1RankNum = d.Rank
			pfm.Year1RankRatio = rankRatio
			pfm.Year1ProfitRatio = d.Syl
		case "2N":
			pfm.Year2Amplitude = d.Avg
			pfm.Year2KindAvg = d.Hs300
			pfm.Year2RankNum = d.Rank
			pfm.Year2RankRatio = rankRatio
			pfm.Year2ProfitRatio = d.Syl
		case "3N":
			pfm.Year3Amplitude = d.Avg
			pfm.Year3KindAvg = d.Hs300
			pfm.Year3RankNum = d.Rank
			pfm.Year3RankRatio = rankRatio
			pfm.Year3ProfitRatio = d.Syl
		case "5N":
			pfm.Year5Amplitude = d.Avg
			pfm.Year5KindAvg = d.Hs300
			pfm.Year5RankNum = d.Rank
			pfm.Year5RankRatio = rankRatio
			pfm.Year5ProfitRatio = d.Syl
		case "JN":
			pfm.ThisYearAmplitude = d.Avg
			pfm.ThisYearKindAvg = d.Hs300
			pfm.ThisYearRankNum = d.Rank
			pfm.ThisYearRankRatio = rankRatio
			pfm.ThisYearProfitRatio = d.Syl
		case "LN":
			pfm.HistoricalAmplitude = d.Avg
			pfm.HistoricalKindAvg = d.Hs300
			pfm.HistoricalRankNum = d.Rank
			pfm.HistoricalRankRatio = rankRatio
			pfm.HistoricalProfitRatio = d.Syl
		}
	}
	fund.Performance = pfm

	// 持仓股票
	stocks := []FundStock{}
	for _, s := range efund.Jjcc.Datas.InverstPosition.FundStocks {
		adj := s.Pctnvchg
		if s.Pctnvchgtype == "减持" {
			adj *= -1
		}
		stock := FundStock{
			Code:        s.Gpdm,
			Name:        s.Gpjc,
			Industry:    s.Indexname,
			HoldRatio:   s.Jzbl,
			AdjustRatio: adj,
		}
		stocks = append(stocks, stock)
	}
	fund.Stocks = stocks

	// 基金经理
	manager := FundManager{}
	if len(efund.Jjjlnew.Datas) == 0 {
		return fund, errors.New("jjjlnew no data")
	}
	jjjl := efund.Jjjlnew.Datas[0]
	if len(jjjl.Manger) == 0 {
		return fund, errors.New("manager no data")
	}
	m := jjjl.Manger[0]
	manager.Name = m.Mgrname
	manager.WorkingDays = m.Totaldays
	manager.ManageDays = jjjl.Days
	manager.ManageRepay = jjjl.Penavgrowth
	manager.YearsAvgRepay = m.Yieldse
	fund.Manager = manager

	// 分红送配
	dividends := []FundDividend{}
	for _, d := range efund.Fhsp.Datas.Fhinfo {
		fd := FundDividend{
			RegDate:    d.Djr,
			Value:      d.Fhfcz,
			RationDate: d.Ffr,
		}
		dividends = append(dividends, fd)
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
		ip := FundIndustryProportion{
			PubDate:  date,
			Industry: []string{},
			Props:    []string{},
		}
		for _, i := range vlist {
			ip.Industry = append(ip.Industry, i["HYMC"])
			ip.Props = append(ip.Props, i["ZJZBL"])
		}
		fund.IndustryProportion = ip
	}

	return fund, nil
}
