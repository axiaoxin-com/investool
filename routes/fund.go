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
	PageNum   int    `json:"page_num"   form:"page_num"`
	PageSize  int    `json:"page_size"  form:"page_size"`
	SortField string `json:"sort_field" form:"sort_field"`
	SortType  string `json:"sort_type"  form:"sort_type"`
}

// FundIndex godoc
func FundIndex(c *gin.Context) {
	p := ParamFundIndex{
		PageNum:   1,
		PageSize:  10,
		SortField: "week_profit",
		SortType:  "asc",
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
	pagedList := []models.FundList{}
	totalCount := len(services.Fund4433List)
	for i := 0; i < totalCount; i += p.PageSize {
		end := i + p.PageSize
		if end > totalCount {
			end = totalCount
		}
		pagedList = append(pagedList, services.Fund4433List[i:end])
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
		"Env":          viper.GetString("env"),
		"Version":      version.Version,
		"PageTitle":    "X-STOCK | 基金",
		"Fund4433List": result,
		"Pagination":   pagi,
		"Param":        p,
	}
	c.HTML(http.StatusOK, "fund_index.html", data)
	return
}
