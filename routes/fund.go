// 首页

package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// FundIndex godoc
func FundIndex(c *gin.Context) {
	data := gin.H{
		"PageTitle": "首页",
	}
	c.HTML(http.StatusOK, "fund.html", data)
	return
}
