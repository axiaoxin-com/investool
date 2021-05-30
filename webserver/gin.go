package webserver

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/axiaoxin-com/goutils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/json-iterator/go/extra"
	"github.com/spf13/viper"
)

func init() {
	// 替换 gin 默认的 validator，更加友好的错误信息
	binding.Validator = &goutils.GinStructValidator{}
	// causes the json binding Decoder to unmarshal a number into an interface{} as a Number instead of as a float64.
	binding.EnableDecoderUseNumber = true

	// jsoniter 启动模糊模式来支持 PHP 传递过来的 JSON。容忍字符串和数字互转
	extra.RegisterFuzzyDecoders()
	// jsoniter 设置支持 private 的 field
	extra.SupportPrivateFields()
}

// NewGinEngine 根据参数创建 gin 的 router engine
// middlewares 需要使用到的中间件列表，默认不为 engine 添加任何中间件
func NewGinEngine(staticsFiles *embed.FS, middlewares ...gin.HandlerFunc) *gin.Engine {
	// set gin mode
	gin.SetMode(viper.GetString("server.mode"))

	engine := gin.New()
	// ///a///b -> /a/b
	engine.RemoveExtraSlash = true

	// use middlewares
	for _, middleware := range middlewares {
		engine.Use(middleware)
	}

	// set template funcmap, must befor load templates
	engine.SetFuncMap(TemplFuncs)

	// load html template
	if staticsFiles != nil {
		tmplPath := viper.GetString("statics.tmplpath")
		if tmplPath != "" {
			t, err := template.ParseFS(*staticsFiles, tmplPath)
			if err != nil {
				panic(err)
			}
			engine.SetHTMLTemplate(t)
		}
		// register statics
		staticsURL := viper.GetString("statics.url")
		if staticsURL != "" {
			engine.StaticFS(staticsURL, http.FS(*staticsFiles))
		}
	}

	return engine
}
