// 学习资料页面

package routes

import (
	"io/ioutil"
	"net/http"

	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/x-stock/version"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
)

// MaterialItem 学习资料具体信息
type MaterialItem struct {
	Name        string `json:"name"`
	DownloadURL string `json:"download_url"`
	Desc        string `json:"desc"`
}

// MaterialSeries 某一个系列的资料
// {
//     "飙股在线等": [
//         MaterialItem, ...
//     ]
// }
type MaterialSeries map[string][]MaterialItem

// TypedMaterialSeries 对MaterialSeries进行分类，如：视频、电子书等
// {
//     "videos": [
//         MaterialSeries, ...
//     ],
//     "ebooks": [
//         MaterialSeries, ...
//     ]
// }
type TypedMaterialSeries map[string][]MaterialSeries

// AllMaterialsList 包含全部资料信息的大JSON列表
// [
//     TypedMaterialSeries, ...
// ]
type AllMaterialsList []TypedMaterialSeries

// MaterialsFilename 资料JSON文件路径
var MaterialsFilename = "./materials.json"

// Materials godoc
func Materials(c *gin.Context) {
	data := gin.H{
		"Env":       viper.GetString("env"),
		"Version":   version.Version,
		"PageTitle": "X-STOCK | 资料",
	}
	f, err := ioutil.ReadFile(MaterialsFilename)
	if err != nil {
		logging.Errorf(c, "Read MaterialsFilename:%v err:%v", MaterialsFilename, err)
		data["Error"] = err
		c.HTML(http.StatusOK, "materials.html", data)
		return
	}
	var mlist AllMaterialsList
	if err := jsoniter.Unmarshal(f, &mlist); err != nil {
		logging.Errorf(c, "json Unmarshal AllMaterialsList err:%v", err)
		data["Error"] = err
		c.HTML(http.StatusOK, "materials.html", data)
		return
	}
	data["AllMaterialsList"] = mlist
	c.HTML(http.StatusOK, "materials.html", data)
	return
}
