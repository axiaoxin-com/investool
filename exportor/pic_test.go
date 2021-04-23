package exportor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExportPic(t *testing.T) {
	e := Exportor{
		Stocks: []Data{
			{
				Name: "中文名称",
				Code: "1234code",
			}, {
				Name: "中文名称1",
				Code: "code12345",
			},
		},
	}

	_, err := e.ExportPic(_ctx, "/tmp/test.png")
	require.Nil(t, err)
}
