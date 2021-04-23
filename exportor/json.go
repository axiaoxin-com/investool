// 导出 json 文件

package exportor

import (
	"context"
	"encoding/json"
	"io/ioutil"
)

// ExportJSON 数据导出为 JSON 文件
func (e Exportor) ExportJSON(ctx context.Context, filename string) (result []byte, err error) {
	result, err = json.MarshalIndent(e, "", "  ")
	err = ioutil.WriteFile(filename, result, 0666)
	return
}
