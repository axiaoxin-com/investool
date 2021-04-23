package exportor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExportExcel(t *testing.T) {
	dummyData := Data{}
	e := Exportor{
		Stocks: []Data{dummyData},
	}
	_, err := e.ExportExcel(_ctx, "/tmp/test.xlsx")
	require.Nil(t, err)
}
