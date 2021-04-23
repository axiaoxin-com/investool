// 导出 csv

package exportor

import (
	"context"
	"os"

	"github.com/gocarina/gocsv"
)

// ExportCSV 数据导出为 CSV 文件
func (e Exportor) ExportCSV(ctx context.Context, filename string) (result []byte, err error) {
	f, err := os.Create(filename)
	if err != nil {
		return
	}
	defer f.Close()
	result, err = gocsv.MarshalBytes(&e)
	err = gocsv.MarshalFile(&e, f)
	return
}
