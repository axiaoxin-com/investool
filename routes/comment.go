// 评论留言

package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Comment godoc
func Comment(c *gin.Context) {
	data := gin.H{
		"Env":       viper.GetString("env"),
		"PageTitle": "X-STOCK | 留言",
	}
	c.HTML(http.StatusOK, "comment.html", data)
	return
}
