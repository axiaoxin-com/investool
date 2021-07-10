// 基金

package routes

import (
	"net/http"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/x-stock/core"
	"github.com/axiaoxin-com/x-stock/models"
	"github.com/axiaoxin-com/x-stock/services"
	"github.com/axiaoxin-com/x-stock/version"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// ParamFundIndex FundIndex 请求参数
type ParamFundIndex struct {
	PageNum  int    `json:"page_num"  form:"page_num"`
	PageSize int    `json:"page_size" form:"page_size"`
	Sort     int    `json:"sort"      form:"sort"`
	Type     string `json:"type"      form:"type"`
}

// FundIndex godoc
func FundIndex(c *gin.Context) {
	fundList := services.Fund4433List
	p := ParamFundIndex{
		PageNum:  1,
		PageSize: 10,
		Sort:     models.FundSortTypeWeek,
	}
	if err := c.ShouldBind(&p); err != nil {
		data := gin.H{
			"Env":       viper.GetString("env"),
			"Version":   version.Version,
			"PageTitle": "X-STOCK | 基金",
			"Error":     err.Error(),
		}
		c.HTML(http.StatusOK, "fund_index.html", data)
		return
	}

	// 过滤
	if p.Type != "" {
		fundList = fundList.FilterByType(p.Type)
	}
	// 排序
	if p.Sort > 0 {
		fundList.Sort(models.FundSortType(p.Sort))
	}
	// 分页
	totalCount := len(fundList)
	pagi := goutils.PaginateByPageNumSize(totalCount, p.PageNum, p.PageSize)
	result := fundList[pagi.StartIndex:pagi.EndIndex]
	data := gin.H{
		"Env":           viper.GetString("env"),
		"Version":       version.Version,
		"PageTitle":     "X-STOCK | 基金 | 4433法则选基",
		"URLPath":       "/fund",
		"FundList":      result,
		"Pagination":    pagi,
		"IndexParam":    p,
		"UpdatedAt":     services.SyncFundTime.Format("2006-01-02 15:04:05"),
		"AllFundCount":  len(services.FundAllList),
		"Fund4433Count": totalCount,
		"FundTypes":     services.Fund4433TypeList,
	}
	c.HTML(http.StatusOK, "fund_index.html", data)
	return
}

// ParamFundFilter FundFilter 请求参数
type ParamFundFilter struct {
	ParamFunderFilter core.ParamFunderFilter
	ParamFundIndex    ParamFundIndex
}

// FundFilter godoc
func FundFilter(c *gin.Context) {
	p := ParamFundFilter{
		ParamFunderFilter: core.ParamFunderFilter{
			MinScale:             2.0,
			MaxScale:             50.0,
			MinManagerYears:      5.0,
			Year1RankRatio:       25.0,
			ThisYear235RankRatio: 25.0,
			Month6RankRatio:      33.33,
			Month3RankRatio:      33.33,
		},
		ParamFundIndex: ParamFundIndex{
			PageNum:  1,
			PageSize: 10,
			Sort:     0,
		},
	}
	if err := c.ShouldBind(&p); err != nil {
		data := gin.H{
			"Env":       viper.GetString("env"),
			"Version":   version.Version,
			"PageTitle": "X-STOCK | 基金 | 基金严选",
			"Error":     err.Error(),
		}
		c.HTML(http.StatusOK, "fund_index.html", data)
		return
	}
	funder := core.NewFunder()
	fundList := funder.Filter(c, p.ParamFunderFilter)
	fundTypes := fundList.Types()
	// 过滤
	if p.ParamFundIndex.Type != "" {
		fundList = fundList.FilterByType(p.ParamFundIndex.Type)
	}
	// 排序
	if p.ParamFundIndex.Sort > 0 {
		fundList.Sort(models.FundSortType(p.ParamFundIndex.Sort))
	}
	// 分页
	pagi := goutils.PaginateByPageNumSize(len(fundList), p.ParamFundIndex.PageNum, p.ParamFundIndex.PageSize)
	result := fundList[pagi.StartIndex:pagi.EndIndex]
	data := gin.H{
		"Env":         viper.GetString("env"),
		"Version":     version.Version,
		"PageTitle":   "X-STOCK | 基金 | 基金严选",
		"URLPath":     "/fund/filter",
		"FundList":    result,
		"Pagination":  pagi,
		"IndexParam":  p.ParamFundIndex,
		"FilterParam": p.ParamFunderFilter,
		"FundTypes":   fundTypes,
	}
	c.HTML(http.StatusOK, "fund_filter.html", data)
	return
}
