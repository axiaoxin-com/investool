// filter 对给定股票进行分析，筛出其中的优质公司。（好公司，但不代表当前股价在涨）
// 选股规则：
// 行业要分散
// 最新 ROE 高于 8%
// ROE 平均值小于 20 时，至少 3 年内逐年递增
// EPS 至少 3 年内逐年递增
// 营业总收入至少 3 年内逐年递增
// 净利润至少 3 年内逐年递增
// 估值较低或中等
// 股价低于合理价格
// 负债率低于 60%

package parser

import (
	"context"
	"fmt"
	"math"
	"sync"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/x-stock/datacenter"
	"github.com/axiaoxin-com/x-stock/datacenter/eastmoney"
	"github.com/axiaoxin-com/x-stock/model"
	"go.uber.org/zap"
)

// MaxWorkerCount 最大并发请求 worker 数
var MaxWorkerCount = 64

// FilterOptions 过滤条件选项
type FilterOptions struct {
	// eastmoney 中的过滤条件
	eastmoney.Filter
	// 连续增长年数
	CheckYears int `json:"check_years"`
	// ROE 高于该值时不做连续增长检查
	NoCheckYearsROE float64 `json:"no_check_years_roe"`
	// 最大负债率百分比(%)
	MaxDebtRatio float64 `json:"max_debt_ratio"`
	// 最小历史波动率
	MinHV float64 `json:"min_hv"`
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
		MinListingYieldYear:      0.0,
		MinPBNewMRQ:              0.0,
		MinPredictNetprofitRatio: 0.0,
		MinPredictIncomeRatio:    0.0,
		MinTotalMarketCap:        0.0,
		Industry:                 "",
		MinPrice:                 0.0,
		MaxPrice:                 0.0,
		ListingOver5Y:            false,
		ExcludeCYB:               true,
		ExcludeKCB:               true,
	},
	CheckYears:      3,
	NoCheckYearsROE: 20,
	MaxDebtRatio:    60,
	MinHV:           0,
}

// StockChecker 判断给定股票是否是好股票
func StockChecker(ctx context.Context, stock model.Stock, options FilterOptions) (defects []string) {
	// 最新 ROE 高于 n%
	if stock.BaseInfo.RoeWeight < options.MinROE {
		defect := fmt.Sprintf(
			"Latest ROE:%v is not greater than:%+v",
			stock.BaseInfo.RoeWeight,
			options.MinROE,
		)
		defects = append(defects, defect)
	}

	// ROE 均值小于 NoCheckYearsROE 时，至少 n 年内逐年递增
	roeavg, err := goutils.AvgFloat64(stock.HistoricalFinaMainData.ROEList(ctx, options.CheckYears))
	if err != nil {
		logging.Warn(ctx, "roe avg error:"+err.Error())
	}
	if roeavg < options.NoCheckYearsROE && !stock.HistoricalFinaMainData.IsIncreasingByYears(ctx, "ROE", options.CheckYears) {
		defect := fmt.Sprintf(
			"ROE is not increasing in %d years. fina:%+v",
			options.CheckYears,
			stock.HistoricalFinaMainData.ROEList(ctx, options.CheckYears),
		)
		defects = append(defects, defect)
	}

	// EPS 至少 n 年内逐年递增
	if !stock.HistoricalFinaMainData.IsIncreasingByYears(ctx, "EPS", options.CheckYears) {
		defect := fmt.Sprintf(
			"EPS is not increasing in %d years. fina:%+v",
			options.CheckYears,
			stock.HistoricalFinaMainData.EPSList(ctx, options.CheckYears),
		)
		defects = append(defects, defect)
	}

	// 营业总收入至少 n 年内逐年递增
	if !stock.HistoricalFinaMainData.IsIncreasingByYears(ctx, "REVENUE", options.CheckYears) {
		defect := fmt.Sprintf(
			"REVENUE is not increasing in %d years. fina:%+v",
			options.CheckYears,
			stock.HistoricalFinaMainData.RevenueList(ctx, options.CheckYears),
		)
		defects = append(defects, defect)
	}

	// 净利润至少 n 年内逐年递增
	if !stock.HistoricalFinaMainData.IsIncreasingByYears(ctx, "PROFIT", options.CheckYears) {
		defect := fmt.Sprintf(
			"PROFIT is not increasing in %d years. fina:%+v",
			options.CheckYears,
			stock.HistoricalFinaMainData.ProfitList(ctx, options.CheckYears),
		)
		defects = append(defects, defect)
	}

	// 估值较低或中等
	if stock.ValuationStatus == eastmoney.ValuationHigh {
		defect := "ValuationStatus is high"
		defects = append(defects, defect)
	}

	// 股价低于合理价格
	if stock.RightPrice != -1 && stock.BaseInfo.NewPrice > stock.RightPrice {
		defect := fmt.Sprintf(
			"NewPrice:%f is higher than RightPrice:%f",
			stock.BaseInfo.NewPrice,
			stock.RightPrice,
		)
		defects = append(defects, defect)
	}

	// 负债率低于 MaxDebtRatio
	if stock.HistoricalFinaMainData[0].Zcfzl > options.MaxDebtRatio {
		defect := fmt.Sprintf(
			"DebtRatio(Zcfzl):%f is greater than %f",
			stock.HistoricalFinaMainData[0].Zcfzl,
			options.MaxDebtRatio,
		)
		defects = append(defects, defect)
	}

	// 历史波动率 （可选条件）
	if options.MinHV != 0 {
		if stock.HistoricalVolatility > options.MinHV {
			defect := fmt.Sprintf(
				"HistoricalVolatility:%f is greater than %f",
				stock.HistoricalVolatility,
				options.MinHV,
			)
			defects = append(defects, defect)
		}
	}
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
			// 按条件判断是否为优质股票
			stock, err := model.NewStock(ctx, baseInfo, false)
			if err != nil {
				logging.Error(ctx, "NewStock error:"+err.Error())
				return
			}

			if defects := StockChecker(ctx, stock, options); len(defects) == 0 {
				result = append(result, stock)
			} else {
				logging.Info(ctx, fmt.Sprintf("%s %s has some defects", stock.BaseInfo.SecurityNameAbbr, stock.BaseInfo.Secucode), zap.Any("defects", defects))
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
