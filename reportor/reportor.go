// Package reportor excel 报表生成器
package reportor

import (
	"context"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

// GenSelectedStocksReport 生成初筛结果报表
func GenSelectedStocksReport(ctx context.Context, filename string) error {
	f := excelize.NewFile()
	index := f.NewSheet("初筛结果")
	f.SetActiveSheet(index)
	if err := f.SaveAs(filename); err != nil {
		return err
	}
	return nil
}
