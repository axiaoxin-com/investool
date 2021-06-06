// 首页

package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Index 返回首页
func Index(c *gin.Context) {
	data := gin.H{}
	c.HTML(http.StatusOK, "files/html/index.html", data)
	return
}
