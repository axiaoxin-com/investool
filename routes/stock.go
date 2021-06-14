// 首页

package routes

import (
	"net/http"

	"github.com/axiaoxin-com/x-stock/core"
	"github.com/axiaoxin-com/x-stock/datacenter/eastmoney"
	"github.com/axiaoxin-com/x-stock/models"
	"github.com/axiaoxin-com/x-stock/services"
	"github.com/gin-gonic/gin"
)

// StockIndex 股票页面
func StockIndex(c *gin.Context) {
	data := gin.H{
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
		"PageTitle": "X-STOCK | 股票 | 基本面筛选",
		"Error":     "",
		"Stocks":    models.StockList{},
	}

	param := ParamStockSelector{}
	if err := c.ShouldBind(&param); err != nil {
		data["Error"] = err.Error()
		c.HTML(http.StatusOK, "stock_index.html", data)
		return
	}
	var checker core.Checker
	if param.FilterWithChecker {
		checker = core.NewChecker(c, param.CheckerOptions)
	}

	selector := core.NewSelector(c, param.Filter, &checker)
	stocks, err := selector.AutoFilterStocks(c)
	if err != nil {
		data["Error"] = err.Error()
		c.HTML(http.StatusOK, "stock_index.html", data)
		return
	}
	dlist := models.ExportorDataList{}
	for _, s := range stocks {
		dlist = append(dlist, models.NewExportorData(c, s))
	}
	data["Stocks"] = dlist
	c.JSON(http.StatusOK, data)
	return
}

// StockChecker 处理个股检测
func StockChecker(c *gin.Context) {
	data := gin.H{
		"PageTitle": "X-STOCK | 股票 | 个股检测",
		"Error":     "",
	}
	c.HTML(http.StatusOK, "stock_index.html", data)
	return
}
