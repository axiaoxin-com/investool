// 关于

package routes

import (
	"net/http"

	"github.com/axiaoxin-com/x-stock/version"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// About godoc
func About(c *gin.Context) {
	data := gin.H{
		"Env":       viper.GetString("env"),
		"Version":   version.Version,
		"PageTitle": "X-STOCK | 关于",
	}
	c.HTML(http.StatusOK, "about.html", data)
	return
}
