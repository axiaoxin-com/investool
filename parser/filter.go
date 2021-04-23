// filter 对给定股票进行分析，筛出其中的优质公司。（好公司，但不代表当前股价在涨）
// 选股规则：
// 0. 行业要分散
// 1. 最新 ROE 高于 8%
// 2. ROE 平均值小于 20 时，至少 3 年内逐年递增
// 3. EPS 至少 3 年内逐年递增
// 4. 营业总收入至少 3 年内逐年递增
// 5. 净利润至少 3 年内逐年递增
// 6. 估值较低或中等
// 7. 股价低于合理价格
// 8. 历史波动率在 1 以内（持仓占比： 0.1:0.1-0.5:0.5-1 = 3:3:4 ）
// 9. 负债率低于 60%

package parser

import (
	"context"
	"math"
	"sync"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/x-stock/datacenter"
	"github.com/axiaoxin-com/x-stock/datacenter/eastmoney"
	"github.com/axiaoxin-com/x-stock/model"
)

// MaxWorkerCount 最大并发请求 worker 数
var MaxWorkerCount = 64

// FilterOptions 过滤条件选项
type FilterOptions struct {
	// eastmoney 中的过滤条件
	eastmoney.Filter
	// 连续增长年数
	CheckYears int
	// ROE 高于该值时不做连续增长检查
	NoCheckYearsROE float64
	// 最大负债率百分比(%)
	MaxDebtRatio float64
	// 最大历史波动率
	MaxHV float64
}

// DefaultFilterOptions 默认过滤条件值
var DefaultFilterOptions = FilterOptions{
	Filter: eastmoney.Filter{
		MinROE:                   8.0,
		MinNetprofitYoyRatio:     0.0,
		MinToiYoyRatio:           0.0,
		MinZXGXL:                 0.0,
		MinNetprofitGrowthrate3Y: 0.0,
		MinIncomeGrowthrate3Y:    0.0,
		MinPredictNetprofitRatio: 0.0,
		MinPredictIncomeRatio:    0.0,
		MinListingYieldYear:      0.0,
		MinPBNewMRQ:              0.0,
		MinTotalMarketCap:        0.0,
		Industry:                 "",
		MinPrice:                 0.0,
		MaxPrice:                 0.0,
		ListingOver5Y:            false,
	},
	CheckYears:      3,
	NoCheckYearsROE: 20,
	MaxDebtRatio:    60,
	MaxHV:           6,
}

// IsGoodStock 判断给定股票是否是好股票
func IsGoodStock(ctx context.Context, baseInfo eastmoney.StockInfo, options FilterOptions) (stock model.Stock, ok bool) {
	stock, err := model.NewStock(ctx, baseInfo)
	if err != nil {
		logging.Error(ctx, "NewStock error:"+err.Error())
		return
	}
	// 1. 最新 ROE 高于 n%
	if stock.BaseInfo.RoeWeight < options.MinROE {
		logging.Debugf(
			ctx,
			"IsGoodStock skip %s %s: Latest ROE:%v is not greater than:%+v",
			stock.BaseInfo.SecurityNameAbbr,
			stock.BaseInfo.Secucode,
			stock.BaseInfo.RoeWeight,
			options.MinROE,
		)
		return
	}

	// 2. ROE 均值小于 20 时，至少 n 年内逐年递增
	roeavg, err := goutils.AvgFloat64(stock.HistoricalFinaMainData.ROEList(ctx, options.CheckYears))
	if err != nil {
		logging.Error(ctx, "roe avg error:"+err.Error())
	}
	if roeavg < options.NoCheckYearsROE && !stock.HistoricalFinaMainData.IsIncreasingByYears(ctx, "ROE", options.CheckYears) {
		logging.Debugf(
			ctx,
			"IsGoodStock skip %s %s: ROE is not increasing in %d years. fina:%+v",
			stock.BaseInfo.SecurityNameAbbr,
			stock.BaseInfo.Secucode,
			options.CheckYears,
			stock.HistoricalFinaMainData.ROEList(ctx, options.CheckYears),
		)
		return
	}

	// 3. EPS 至少 n 年内逐年递增
	if !stock.HistoricalFinaMainData.IsIncreasingByYears(ctx, "EPS", options.CheckYears) {
		logging.Debugf(
			ctx,
			"IsGoodStock skip %s %s EPS is not increasing in %d years. fina:%+v",
			stock.BaseInfo.SecurityNameAbbr,
			stock.BaseInfo.Secucode,
			options.CheckYears,
			stock.HistoricalFinaMainData.EPSList(ctx, options.CheckYears),
		)
		return
	}

	// 4. 营业总收入至少 n 年内逐年递增
	if !stock.HistoricalFinaMainData.IsIncreasingByYears(ctx, "REVENUE", options.CheckYears) {
		logging.Debugf(
			ctx,
			"IsGoodStock skip %s %s REVENUE is not increasing in %d years. fina:%+v",
			stock.BaseInfo.SecurityNameAbbr,
			stock.BaseInfo.Secucode,
			options.CheckYears,
			stock.HistoricalFinaMainData.RevenueList(ctx, options.CheckYears),
		)
		return
	}

	// 5. 净利润至少 n 年内逐年递增
	if !stock.HistoricalFinaMainData.IsIncreasingByYears(ctx, "PROFIT", options.CheckYears) {
		logging.Debugf(
			ctx,
			"IsGoodStock skip %s %s PROFIT is not increasing in %d years. fina:%+v",
			stock.BaseInfo.SecurityNameAbbr,
			stock.BaseInfo.Secucode,
			options.CheckYears,
			stock.HistoricalFinaMainData.ProfitList(ctx, options.CheckYears),
		)
		return
	}

	// 6. 估值较低或中等
	if stock.ValuationStatus == eastmoney.ValuationHigh {
		logging.Debugf(
			ctx,
			"IsGoodStock skip %s %s ValuationStatus is high",
			stock.BaseInfo.SecurityNameAbbr,
			stock.BaseInfo.Secucode,
		)
		return
	}

	// 7. 股价低于合理价格
	if stock.RightPrice != -1 && stock.BaseInfo.NewPrice > stock.RightPrice {
		logging.Debugf(
			ctx,
			"IsGoodStock skip %s %s NewPrice:%f is higher than RightPrice:%f",
			stock.BaseInfo.SecurityNameAbbr,
			stock.BaseInfo.Secucode,
			stock.BaseInfo.NewPrice,
			stock.RightPrice,
		)
		return
	}

	// 8. 历史波动率在 1 以内
	if stock.HistoricalVolatility > options.MaxHV {
		logging.Debugf(
			ctx,
			"IsGoodStock skip %s %s HistoricalVolatility:%f is greater than 1",
			stock.BaseInfo.SecurityNameAbbr,
			stock.BaseInfo.Secucode,
			stock.HistoricalVolatility,
		)
		return
	}

	// 9. 负债率低于 60%
	if stock.HistoricalFinaMainData[0].Zcfzl > options.MaxDebtRatio {
		logging.Debugf(
			ctx,
			"IsGoodStock skip %s %s Zcfzl:%f is greater than %f",
			stock.BaseInfo.SecurityNameAbbr,
			stock.BaseInfo.Secucode,
			stock.HistoricalFinaMainData[0].Zcfzl,
			options.MaxDebtRatio,
		)
		return
	}
	ok = true
	return
}

// AutoFilterStocks 按默认设置自动筛选股票
func AutoFilterStocks(ctx context.Context) (model.StockList, error) {
	return AutoFilterStocksWithOptions(ctx, DefaultFilterOptions)
}

// AutoFilterStocksWithOptions 按设置自动筛选股票
func AutoFilterStocksWithOptions(ctx context.Context, options FilterOptions) (result model.StockList, err error) {
	stocks, err := datacenter.EastMoney.QuerySelectedStocksWithFilter(ctx, options.Filter)
	if err != nil {
		return
	}
	logging.Infof(ctx, "AutoFilterStock will filter from %d stocks", len(stocks))

	// 最多 MaxWorkerCount 个 groutine 并发执行筛选任务
	workerCount := int(math.Min(float64(len(stocks)), float64(MaxWorkerCount)))
	jobChan := make(chan struct{}, workerCount)
	wg := sync.WaitGroup{}

	for _, baseInfo := range stocks {
		wg.Add(1)
		jobChan <- struct{}{}

		go func(ctx context.Context, baseInfo eastmoney.StockInfo) {
			defer func() {
				wg.Done()
				<-jobChan
				if r := recover(); r != nil {
					logging.Errorf(ctx, "recover from:%v", r)
				}
			}()
			// 按添加判断是否为优质股票
			if stock, ok := IsGoodStock(ctx, baseInfo, options); ok {
				result = append(result, stock)
			}
		}(ctx, baseInfo)
	}
	wg.Wait()
	logging.Infof(ctx, "AutoFilterStock selected %d stocks", len(result))
	result.SortByROE()
	return
}

// AutoFilterStocksByIndustry 自动按行业选择优质股票
func AutoFilterStocksByIndustry(ctx context.Context, options FilterOptions) (result model.StockList, err error) {
	return
}
