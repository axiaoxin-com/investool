// 首页

package routes

import (
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/x-stock/core"
	"github.com/axiaoxin-com/x-stock/datacenter/eastmoney"
	"github.com/axiaoxin-com/x-stock/models"
	"github.com/axiaoxin-com/x-stock/services"
	"github.com/axiaoxin-com/x-stock/version"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// StockIndex 股票页面
func StockIndex(c *gin.Context) {
	data := gin.H{
		"Env":          viper.GetString("env"),
		"Version":      version.Version,
		"PageTitle":    "X-STOCK | 股票",
		"Error":        "",
		"IndustryList": services.StockIndustryList,
	}
	c.HTML(http.StatusOK, "stock_index.html", data)
	return
}

// ParamStockSelector StockSelector 请求参数
type ParamStockSelector struct {
	Filter            eastmoney.Filter
	CheckerOptions    core.CheckerOptions
	FilterWithChecker bool `form:"selector_with_checker"`
}

// StockSelector 返回基本面筛选结果json
func StockSelector(c *gin.Context) {
	data := gin.H{
		"Env":       viper.GetString("env"),
		"Version":   version.Version,
		"PageTitle": "X-STOCK | 股票 | 基本面筛选",
		"Error":     "",
		"Stocks":    models.StockList{},
	}

	param := ParamStockSelector{}
	if err := c.ShouldBind(&param); err != nil {
		data["Error"] = err.Error()
		c.JSON(http.StatusOK, data)
		return
	}
	var checker *core.Checker
	if param.FilterWithChecker {
		checker = core.NewChecker(c, param.CheckerOptions)
	}

	selector := core.NewSelector(c, param.Filter, checker)
	stocks, err := selector.AutoFilterStocks(c)
	if err != nil {
		data["Error"] = err.Error()
		c.JSON(http.StatusOK, data)
		return
	}
	stocks.SortByROE()
	dlist := models.ExportorDataList{}
	for _, s := range stocks {
		dlist = append(dlist, models.NewExportorData(c, s))
	}
	data["Stocks"] = dlist
	c.JSON(http.StatusOK, data)
	return
}

// ParamStockChecker StockChecker 请求参数
type ParamStockChecker struct {
	Keyword        string `form:"checker_keyword"`
	CheckerOptions core.CheckerOptions
}

// StockChecker 处理个股检测
func StockChecker(c *gin.Context) {
	data := gin.H{
		"Env":       viper.GetString("env"),
		"Version":   version.Version,
		"PageTitle": "X-STOCK | 股票 | 个股检测",
		"Error":     "",
	}
	param := ParamStockChecker{}
	if err := c.ShouldBind(&param); err != nil {
		data["Error"] = err.Error()
		c.JSON(http.StatusOK, data)
		return
	}
	if param.Keyword == "" {
		data["Error"] = "请填写股票代码或简称"
		c.JSON(http.StatusOK, data)
		return
	}
	searcher := core.NewSearcher(c)
	keywords := goutils.SplitStringFields(param.Keyword)
	if len(keywords) > 50 {
		data["Error"] = "股票数量超过限制"
		c.JSON(http.StatusOK, data)
		return
	}
	stocks, err := searcher.SearchStocks(c, keywords)
	if err != nil {
		data["Error"] = err.Error()
		c.JSON(http.StatusOK, data)
		return
	}
	checker := core.NewChecker(c, param.CheckerOptions)
	results := []core.CheckResult{}
	stockNames := []string{}
	finaReportNames := []string{}
	finaAppointPublishDates := []string{}
	lines := [][]gin.H{}
	mainInflows := []string{}

	type dataset struct {
		Label string    `json:"label"`
		Data  []float64 `json:"data"`
	}
	for _, stock := range stocks {
		result, _ := checker.CheckFundamentals(c, stock)
		results = append(results, result)
		stockName := fmt.Sprintf("%s-%s", stock.BaseInfo.SecurityNameAbbr, stock.BaseInfo.Secucode)
		stockNames = append(stockNames, stockName)
		mainInflows = append(mainInflows, stock.MainMoneyNetInflows.String())

		finaReportName := ""
		if len(stock.HistoricalFinaMainData) > 0 {
			finaReportName = stock.HistoricalFinaMainData[0].ReportDateName
		}
		finaReportNames = append(finaReportNames, finaReportName)

		finaAppointPublishDates = append(finaAppointPublishDates, strings.Split(stock.FinaAppointPublishDate, " ")[0])

		roeList := goutils.ReversedFloat64Slice(
			stock.HistoricalFinaMainData.ValueList(
				c,
				eastmoney.ValueListTypeROE,
				param.CheckerOptions.CheckYears,
				eastmoney.FinaReportTypeYear,
			))
		yearLabels := []string{}
		year := time.Now().Year()
		yearcount := int(math.Min(float64(param.CheckerOptions.CheckYears), float64(len(roeList))))
		for i := yearcount; i > 0; i-- {
			yearLabels = append(yearLabels, fmt.Sprint(year-i))
		}
		line0 := gin.H{
			"title":  "",
			"xLable": "年",
			"yLabel": "",
			"data": gin.H{
				"labels": yearLabels,
				"datasets": []dataset{
					{
						Label: "ROE",
						Data:  roeList,
					},
					{
						Label: "EPS",
						Data: goutils.ReversedFloat64Slice(stock.HistoricalFinaMainData.ValueList(
							c,
							eastmoney.ValueListTypeEPS,
							param.CheckerOptions.CheckYears,
							eastmoney.FinaReportTypeYear,
						)),
					},
					{
						Label: "ROA",
						Data: goutils.ReversedFloat64Slice(stock.HistoricalFinaMainData.ValueList(
							c,
							eastmoney.ValueListTypeROA,
							param.CheckerOptions.CheckYears,
							eastmoney.FinaReportTypeYear,
						)),
					},
					{
						Label: "毛利率",
						Data: goutils.ReversedFloat64Slice(stock.HistoricalFinaMainData.ValueList(
							c,
							eastmoney.ValueListTypeMLL,
							param.CheckerOptions.CheckYears,
							eastmoney.FinaReportTypeYear,
						)),
					},
					{
						Label: "净利率",
						Data: goutils.ReversedFloat64Slice(stock.HistoricalFinaMainData.ValueList(
							c,
							eastmoney.ValueListTypeJLL,
							param.CheckerOptions.CheckYears,
							eastmoney.FinaReportTypeYear,
						)),
					},
				},
			},
		}
		line1 := gin.H{
			"title":  "",
			"xLable": "年",
			"yLabel": "",
			"data": gin.H{
				"labels": yearLabels,
				"datasets": []dataset{
					{
						Label: "营收",
						Data: goutils.ReversedFloat64Slice(stock.HistoricalFinaMainData.ValueList(
							c,
							eastmoney.ValueListTypeRevenue,
							param.CheckerOptions.CheckYears,
							eastmoney.FinaReportTypeYear,
						)),
					},
					{
						Label: "毛利",
						Data: goutils.ReversedFloat64Slice(stock.HistoricalFinaMainData.ValueList(
							c,
							eastmoney.ValueListTypeGrossProfit,
							param.CheckerOptions.CheckYears,
							eastmoney.FinaReportTypeYear,
						)),
					},
					{
						Label: "净利",
						Data: goutils.ReversedFloat64Slice(stock.HistoricalFinaMainData.ValueList(
							c,
							eastmoney.ValueListTypeNetProfit,
							param.CheckerOptions.CheckYears,
							eastmoney.FinaReportTypeYear,
						)),
					},
				},
			},
		}
		lines = append(lines, []gin.H{line0, line1})
	}
	data["Results"] = results
	data["StockNames"] = stockNames
	data["FinaReportNames"] = finaReportNames
	data["FinaAppointPublishDates"] = finaAppointPublishDates
	data["Lines"] = lines
	data["MainMoneyNetInflows"] = mainInflows
	c.JSON(http.StatusOK, data)
	return
}
