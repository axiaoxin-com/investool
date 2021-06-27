// 默认实现的 ping api

package routes

import (
	"github.com/axiaoxin-com/x-stock/routes/response"
	"github.com/axiaoxin-com/x-stock/services"
	"github.com/axiaoxin-com/x-stock/version"

	"github.com/gin-gonic/gin"
)

// Ping godoc
// @Summary 默认的 Ping 接口
// @Description 返回 server 相关信息，可以用于健康检查
// @Tags x
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param trace_id header string false "you can set custom trace id in header"
// @Router /x/ping [get]
func Ping(c *gin.Context) {
	data := gin.H{
		"version":             version.Version,
		"mysql":               services.CheckMySQL(c),
		"redis":               services.CheckRedis(c),
		"atomic_level_server": services.CheckAtomicLevelServer(c),
	}
	response.JSON(c, data)
	return
}
