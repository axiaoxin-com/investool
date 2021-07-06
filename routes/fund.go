// 基金

package routes

import (
	"net/http"

	"github.com/axiaoxin-com/goutils"
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

	result := models.FundList{}
	if p.PageSize < 0 {
		result = fundList
		pagi = goutils.PaginateByPageNumSize(totalCount, 1, totalCount)
	} else {
		offset := (p.PageNum - 1) * p.PageSize
		if offset < 0 {
			offset = 0
		}
		if offset <= totalCount {
			end := offset + pagi.PageSize
			if end > totalCount {
				end = totalCount
			}
			result = fundList[offset:end]
		}
	}
	data := gin.H{
		"Env":           viper.GetString("env"),
		"Version":       version.Version,
		"PageTitle":     "X-STOCK | 基金 | 4433法则选基",
		"Fund4433List":  result,
		"Pagination":    pagi,
		"Param":         p,
		"UpdatedAt":     services.SyncFundTime.Format("2006-01-02 15:04:05"),
		"AllFundCount":  len(services.FundAllList),
		"Fund4433Count": totalCount,
		"FundTypes":     services.Fund4433List.Types(),
	}
	c.HTML(http.StatusOK, "fund_index.html", data)
	return
}
