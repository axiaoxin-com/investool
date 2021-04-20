package eastmoney

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueryFinaMainData(t *testing.T) {
	data, err := _em.QueryFinaMainData(_ctx, "002050.SZ")
	require.Nil(t, err)
	require.NotEmpty(t, data)
	data1 := data.FilterByReportType(_ctx, "年报")
	require.NotEmpty(t, data1)
	data2 := data.FilterByReportYear(_ctx, "2020")
	require.Equal(t, 4, len(data2))
}

func TestQueryFinaPublishDate(t *testing.T) {
	date, err := _em.QueryFinaPublishDate(_ctx, "002459")
	require.Nil(t, err)
	t.Log(date)
}
