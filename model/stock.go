// 股票对象封装

package model

import (
	"context"
	"sort"

	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/x-stock/datacenter"
	"github.com/axiaoxin-com/x-stock/datacenter/eastmoney"
	"github.com/axiaoxin-com/x-stock/datacenter/eniu"
)

// Stock 接口返回的股票信息结构
type Stock struct {
	// 东方财富接口返回的基本信息
	BaseInfo eastmoney.StockInfo
	// 历史财报信息
	HistoricalFinaMainData eastmoney.HistoricalFinaMainData `json:"historical_fina_main_data"`
	// 市盈率、市净率、市销率、市现率估值
	ValuationMap map[string]string
	// 4 种估值综合估值状态
	ValuationStatus float64
	// 历史市盈率
	HistoricalPEList eastmoney.HistoricalPEList
	// 合理价格：历史市盈率中位数 * (EPS * (1 + 今年 Q1 营收增长比))
	RightPrice float64
	// 历史股价
	HistoricalPrice eniu.RespHistoricalStockPrice
	// 历史波动率
	HistoricalVolatility float64
	// 公司资料
	CompanyProfile eastmoney.CompanyProfile
	// 预约财报披露日期
	FinaAppointPublishDate string
	// 机构评级
	OrgRatingList eastmoney.OrgRatingList
	// 盈利预测
	ProfitPredictList eastmoney.ProfitPredictList
}

// ValuationStatusDesc 综合估值状态文字描述
func (s Stock) ValuationStatusDesc() string {
	desc := ""
	switch s.ValuationStatus {
	case eastmoney.ValuationLow:
		desc = "估值较低"
	case eastmoney.ValuationModerate:
		desc = "估值中等"
	case eastmoney.ValuationHigh:
		desc = "估值较高"
	}
	return desc
}

// GetPrice 返回股价，没开盘时可能是字符串“-”，此时返回最近历史股价，无历史价则返回 -1
func (s Stock) GetPrice() float64 {
	p, ok := s.BaseInfo.NewPrice.(float64)
	if ok {
		return p
	}
	if len(s.HistoricalPrice.Price) == 0 {
		return -1.0
	}
	return s.HistoricalPrice.Price[len(s.HistoricalPrice.Price)-1]
}

// StockList 股票列表
type StockList []Stock

// SortByROE 股票列表按 ROE 排序
func (s StockList) SortByROE() {
	sort.Slice(s, func(i, j int) bool {
		return s[i].BaseInfo.RoeWeight > s[j].BaseInfo.RoeWeight
	})
}

// NewStock 创建 Stock 对象
func NewStock(ctx context.Context, baseInfo eastmoney.StockInfo, strict bool) (Stock, error) {
	s := Stock{
		BaseInfo: baseInfo,
	}

	// 获取财报
	hf, err := datacenter.EastMoney.QueryHistoricalFinaMainData(ctx, s.BaseInfo.Secucode)
	if err != nil {
		if strict {
			return s, err
		}
		logging.Warn(ctx, err.Error())
	}
	s.HistoricalFinaMainData = hf

	// 获取综合估值
	status, valMap, err := datacenter.EastMoney.QueryValuationStatus(ctx, s.BaseInfo.Secucode)
	if err != nil {
		if strict {
			return s, err
		}
		logging.Warn(ctx, err.Error())
	}
	s.ValuationMap = valMap
	s.ValuationStatus = status

	// 历史市盈率
	peList, err := datacenter.EastMoney.QueryHistoricalPEList(ctx, s.BaseInfo.Secucode)
	if err != nil {
		if strict {
			return s, err
		}
		logging.Warn(ctx, err.Error())
	}
	s.HistoricalPEList = peList

	// 合理价格
	// 今年一季报营收增长比
	ratio, err := s.HistoricalFinaMainData.Q1RevenueIncreasingRatio(ctx)
	if err == nil {
		if strict {
			return s, err
		}
		peMidVal, err := peList.GetMidValue(ctx)
		if err != nil {
			logging.Warn(ctx, err.Error())
			s.RightPrice = -1
		}
		s.RightPrice = peMidVal * (s.HistoricalFinaMainData[0].Epsjb * (1 + ratio))
	} else {
		// 一季报没有发布则设置合理价为 -1
		s.RightPrice = -1
	}

	// 历史股价
	hisPrice, err := datacenter.Eniu.QueryHistoricalStockPrice(ctx, s.BaseInfo.Secucode)
	if err != nil {
		if strict {
			return s, err
		}
		logging.Warn(ctx, err.Error())
	}
	s.HistoricalPrice = hisPrice

	// 历史波动率
	hv, err := hisPrice.HistoricalVolatility(ctx, "YEAR")
	if err != nil {
		if strict {
			return s, err
		}
		logging.Warn(ctx, err.Error())
	}
	s.HistoricalVolatility = hv

	// 公司资料
	cp, err := datacenter.EastMoney.QueryCompanyProfile(ctx, s.BaseInfo.Secucode)
	if err != nil {
		if strict {
			return s, err
		}
		logging.Warn(ctx, err.Error())
	}
	s.CompanyProfile = cp

	// 最新财报预约披露时间
	finaPubDate, err := datacenter.EastMoney.QueryAppointFinaPublishDate(ctx, s.BaseInfo.SecurityCode)
	if err != nil {
		if strict {
			return s, err
		}
		logging.Warn(ctx, err.Error())
	}
	s.FinaAppointPublishDate = finaPubDate

	// 机构评级统计
	orgRatings, err := datacenter.EastMoney.QueryOrgRating(ctx, s.BaseInfo.Secucode)
	if err != nil {
		if strict {
			return s, err
		}
		logging.Warn(ctx, err.Error())
	}
	s.OrgRatingList = orgRatings

	// 盈利预测
	pps, err := datacenter.EastMoney.QueryProfitPredict(ctx, s.BaseInfo.Secucode)
	if err != nil {
		if strict {
			return s, err
		}
		logging.Warn(ctx, err.Error())
	}
	s.ProfitPredictList = pps
	return s, nil
}
