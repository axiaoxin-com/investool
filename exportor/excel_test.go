package exportor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExportExcel(t *testing.T) {
	e := Exportor{}
	_, err := e.ExportExcel(_ctx, "./test.xlsx")
	require.Nil(t, err)
}
