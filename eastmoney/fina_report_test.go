package eastmoney

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueryMainFinaData(t *testing.T) {
	data, err := _em.QueryMainFinaData(_ctx, "002050.SZ")
	require.Nil(t, err)
	require.NotEmpty(t, data)
	data1 := data.FilterByReportType(_ctx, "年报")
	require.NotEmpty(t, data1)
	data2 := data.FilterByReportYear(_ctx, "2020")
	require.Equal(t, 4, len(data2))
}
