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
	PageNum  int `json:"page_num"  form:"page_num"`
	PageSize int `json:"page_size" form:"page_size"`
	Sort     int `json:"sort"      form:"sort"`
}

// FundIndex godoc
func FundIndex(c *gin.Context) {
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
	totalCount := len(services.Fund4433List)
	pagedList := []models.FundList{}
	if p.PageSize < 0 {
		pagedList = append(pagedList, services.Fund4433List)
		totalCount = 1
	} else {
		for i := 0; i < totalCount; i += p.PageSize {
			end := i + p.PageSize
			if end > totalCount {
				end = totalCount
			}
			pagedList = append(pagedList, services.Fund4433List[i:end])
		}
	}

	result := models.FundList{}
	index := 0
	if len(pagedList) > 0 {
		if p.PageNum < 1 || p.PageNum > len(pagedList) {
			index = 0
		} else {
			index = p.PageNum - 1
		}
		result = pagedList[index]
	}
	pagi := goutils.PaginateByPageNumSize(totalCount, p.PageNum, p.PageSize)
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
	}
	c.HTML(http.StatusOK, "fund_index.html", data)
	return
}
