// 基金

package routes

import (
	"net/http"

	"github.com/axiaoxin-com/x-stock/version"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// FundIndex godoc
func FundIndex(c *gin.Context) {
	data := gin.H{
		"Env":       viper.GetString("env"),
		"Version":   version.Version,
		"PageTitle": "X-STOCK | 基金",
	}
	c.HTML(http.StatusOK, "fund_index.html", data)
	return
}
