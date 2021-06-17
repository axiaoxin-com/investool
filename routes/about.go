// 关于

package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// About godoc
func About(c *gin.Context) {
	data := gin.H{
		"Env":       viper.GetString("env"),
		"PageTitle": "X-STOCK | 关于",
	}
	c.HTML(http.StatusOK, "about.html", data)
	return
}
