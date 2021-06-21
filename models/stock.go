// 股票对象封装

package models

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
	BaseInfo eastmoney.StockInfo `json:"base_info"`
	// 历史财报信息
	HistoricalFinaMainData eastmoney.HistoricalFinaMainData `json:"historical_fina_main_data"`
	// 市盈率、市净率、市销率、市现率估值
	ValuationMap map[string]string `json:"valuation_map"`
	// 历史市盈率
	HistoricalPEList eastmoney.HistoricalPEList `json:"historical_pe_list"`
	// 合理价格：历史市盈率中位数 * (去年EPS * (1 + 今年 Q1 营收增长比))
	RightPrice float64 `json:"right_price"`
	// 历史股价
	HistoricalPrice eniu.RespHistoricalStockPrice `json:"historical_price"`
	// 历史波动率
	HistoricalVolatility float64 `json:"historical_volatility"`
	// 公司资料
	CompanyProfile eastmoney.CompanyProfile `json:"company_profile"`
	// 预约财报披露日期
	FinaAppointPublishDate string `json:"fina_appoint_publish_date"`
	// 实际财报披露日期
	FinaActualPublishDate string `json:"fina_actual_publish_date"`
	// 财报披露日期
	FinaReportDate string `json:"fina_report_date"`
	// 机构评级
	OrgRatingList eastmoney.OrgRatingList `json:"org_rating_list"`
	// 盈利预测
	ProfitPredictList eastmoney.ProfitPredictList `json:"profit_predict_list"`
	// 价值评估
	JZPG eastmoney.JZPG `json:"jzpg"`
	// PEG=PE/净利润复合增长率
	PEG float64 `json:"peg"`
	// 历史利润表
	HistoricalGincomeList eastmoney.GincomeDataList `json:"historical_gincome_list"`
	// 本业营收比=营业利润/(营业利润+营业外收入)
	BYYSRatio float64 `json:"byys_ratio"`
	// 最新财报审计意见
	FinaReportOpinion string `json:"fina_report_opinion"`
	// 历史现金流量表
	HistoricalCashflowList eastmoney.CashflowDataList `json:"historical_cashdlow_list"`
	// 最新经营活动产生的现金流量净额
	NetcashOperate float64 `json:"netcash_operate"`
	// 最新投资活动产生的现金流量净额
	NetcashInvest float64 `json:"netcash_invest"`
	// 最新筹资活动产生的现金流量净额
	NetcashFinance float64 `json:"netcash_finance"`
	// 自由现金流
	NetcashFree float64 `json:"netcash_free"`
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

// GetOrgType 获取机构类型
func (s Stock) GetOrgType() string {
	if len(s.HistoricalFinaMainData) == 0 {
		return ""
	}
	return s.HistoricalFinaMainData[0].OrgType
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
		logging.Warn(ctx, "NewStock QueryHistoricalFinaMainData err:"+err.Error())
	}
	s.HistoricalFinaMainData = hf

	// 获取综合估值
	valMap, err := datacenter.EastMoney.QueryValuationStatus(ctx, s.BaseInfo.Secucode)
	if err != nil {
		if strict {
			return s, err
		}
		logging.Warn(ctx, "NewStock QueryValuationStatus err:"+err.Error())
	}
	s.ValuationMap = valMap

	// 历史市盈率
	peList, err := datacenter.EastMoney.QueryHistoricalPEList(ctx, s.BaseInfo.Secucode)
	if err != nil {
		if strict {
			return s, err
		}
		logging.Warn(ctx, "NewStock QueryHistoricalPEList err:"+err.Error())
	}
	s.HistoricalPEList = peList

	// 合理价格判断，一季报没有发布则设置合理价为 -1
	s.RightPrice = -1
	// 今年一季报营收增长比
	ratio, err := s.HistoricalFinaMainData.Q1RevenueIncreasingRatio(ctx)
	if err != nil {
		if strict {
			return s, err
		}
		logging.Warn(ctx, "NewStock Q1RevenueIncreasingRatio err:"+err.Error())
	}
	// pe 中位数
	peMidVal, err := peList.GetMidValue(ctx)
	if err != nil {
		if strict {
			return s, err
		}
		logging.Warn(ctx, "NewStock GetMidValue err:"+err.Error())
	}
	reports := s.HistoricalFinaMainData.FilterByReportType(ctx, "年报")
	if len(reports) > 0 {
		s.RightPrice = peMidVal * (reports[0].Epsjb * (1 + ratio/100.0))
	}

	// 历史股价
	hisPrice, err := datacenter.Eniu.QueryHistoricalStockPrice(ctx, s.BaseInfo.Secucode)
	if err != nil {
		if strict {
			return s, err
		}
		logging.Warn(ctx, "NewStock QueryHistoricalStockPrice err:"+err.Error())
	}
	s.HistoricalPrice = hisPrice

	// 历史波动率
	hv, err := hisPrice.HistoricalVolatility(ctx, "YEAR")
	if err != nil {
		if strict {
			return s, err
		}
		logging.Warn(ctx, "NewStock HistoricalVolatility err:"+err.Error())
	}
	s.HistoricalVolatility = hv

	// 公司资料
	cp, err := datacenter.EastMoney.QueryCompanyProfile(ctx, s.BaseInfo.Secucode)
	if err != nil {
		if strict {
			return s, err
		}
		logging.Warn(ctx, "NewStock QueryCompanyProfile err:"+err.Error())
	}
	s.CompanyProfile = cp

	// 最新财报预约披露时间
	finaPubDateList, err := datacenter.EastMoney.QueryFinaPublishDateList(ctx, s.BaseInfo.SecurityCode)
	if err != nil {
		if strict {
			return s, err
		}
		logging.Warn(ctx, "NewStock QueryFinaPublishDateList err:"+err.Error())
	}
	if len(finaPubDateList) > 0 {
		s.FinaAppointPublishDate = finaPubDateList[0].AppointPublishDate
		s.FinaActualPublishDate = finaPubDateList[0].ActualPublishDate
		s.FinaReportDate = finaPubDateList[0].ReportDate
	}

	// 机构评级统计
	orgRatings, err := datacenter.EastMoney.QueryOrgRating(ctx, s.BaseInfo.Secucode)
	if err != nil {
		if strict {
			return s, err
		}
		logging.Warn(ctx, "NewStock QueryOrgRating err:"+err.Error())
	}
	s.OrgRatingList = orgRatings

	// 盈利预测
	pps, err := datacenter.EastMoney.QueryProfitPredict(ctx, s.BaseInfo.Secucode)
	if err != nil {
		if strict {
			return s, err
		}
		logging.Warn(ctx, "NewStock QueryProfitPredict err:"+err.Error())
	}
	s.ProfitPredictList = pps

	// 价值评估
	jzpg, err := datacenter.EastMoney.QueryJiaZhiPingGu(ctx, s.BaseInfo.Secucode)
	if err != nil {
		if strict {
			return s, err
		}
		logging.Warn(ctx, "NewStock QueryJiaZhiPingGu err:"+err.Error())
	}
	s.JZPG = jzpg

	// PEG
	s.PEG = s.BaseInfo.PE / s.BaseInfo.NetprofitGrowthrate3Y

	// 利润表数据
	gincomeList, err := datacenter.EastMoney.QueryFinaGincomeData(ctx, s.BaseInfo.Secucode)
	if err != nil {
		if strict {
			return s, err
		}
		logging.Warn(ctx, "NewStock QueryFinaGincomeData err:"+err.Error())
	}
	s.HistoricalGincomeList = gincomeList
	if len(s.HistoricalGincomeList) > 0 {
		// 本业营收比
		gincome := s.HistoricalGincomeList[0]
		s.BYYSRatio = gincome.OperateProfit / (gincome.OperateProfit + gincome.NonbusinessIncome)
		// 审计意见
		s.FinaReportOpinion = gincome.OpinionType
	}

	// 现金流量表数据
	cashflow, err := datacenter.EastMoney.QueryFinaCashflowData(ctx, s.BaseInfo.Secucode)
	if err != nil {
		if strict {
			return s, err
		}
		logging.Warn(ctx, "NewStock QueryFinaCashflowData err:"+err.Error())
	}
	s.HistoricalCashflowList = cashflow
	if len(s.HistoricalCashflowList) > 0 {
		cf := s.HistoricalCashflowList[0]
		s.NetcashOperate = cf.NetcashOperate
		s.NetcashInvest = cf.NetcashInvest
		s.NetcashFinance = cf.NetcashFinance
		if cf.NetcashInvest < 0 {
			s.NetcashFree = s.NetcashOperate + s.NetcashInvest
		} else {
			s.NetcashFree = s.NetcashOperate - s.NetcashInvest
		}
	}

	return s, nil
}
