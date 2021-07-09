// web 服务

package cmds

import (
	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/x-stock/cron"
	"github.com/axiaoxin-com/x-stock/routes"
	"github.com/axiaoxin-com/x-stock/routes/response"
	"github.com/axiaoxin-com/x-stock/services"
	"github.com/axiaoxin-com/x-stock/webserver"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

const (
	// ProcessorWebserver web 服务
	ProcessorWebserver = "webserver"
)

// FlagsWebserver cli flags
func FlagsWebserver() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "config",
			Aliases:  []string{"c"},
			Value:    "./config.toml",
			Usage:    "配置文件",
			Required: false,
		},
	}
}

// DefaultGinMiddlewares 默认的 gin server 使用的中间件列表
func DefaultGinMiddlewares() []gin.HandlerFunc {
	m := []gin.HandlerFunc{
		// 记录请求处理日志，最顶层执行
		webserver.GinLogMiddleware(),
		// 捕获 panic 保存到 context 中由 GinLogger 统一打印， panic 时返回 500 JSON
		webserver.GinRecovery(response.Respond),
	}

	// 配置开启请求限频则添加限频中间件
	if viper.GetBool("ratelimiter.enable") {
		m = append(m, webserver.GinRatelimitMiddleware())
	}
	return m
}

// ActionWebserver cli action
func ActionWebserver() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		configFile := c.String("config")
		webserver.InitWithConfigFile(configFile)

		// 初始化或加载外部依赖服务
		if err := services.Init(); err != nil {
			logging.Error(nil, "services init error:"+err.Error())
		}

		// 启动定时任务
		if viper.GetString("env") == "prod" {
			cron.RunCronJobs(true)
		}
		// 创建 gin app
		middlewares := DefaultGinMiddlewares()
		server := webserver.NewGinEngine(middlewares...)
		// 注册路由
		routes.Register(server)
		// 运行服务
		webserver.Run(server)
		return nil
	}
}

// CommandWebserver 检测器 cli command
func CommandWebserver() *cli.Command {
	flags := FlagsWebserver()
	cmd := &cli.Command{
		Name:   ProcessorWebserver,
		Usage:  "web服务器",
		Flags:  flags,
		Action: ActionWebserver(),
	}
	return cmd
}
