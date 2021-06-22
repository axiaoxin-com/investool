package eastmoney

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueryHistoricalFinaMainData(t *testing.T) {
	data, err := _em.QueryHistoricalFinaMainData(_ctx, "600188.SH")
	require.Nil(t, err)
	require.NotEmpty(t, data)
	data1 := data.FilterByReportType(_ctx, FinaReportTypeYear)
	require.NotEmpty(t, data1)
	data2 := data.FilterByReportYear(_ctx, 2020)
	require.Equal(t, 4, len(data2))
	ratio, err := data.Q1RevenueIncreasingRatio(_ctx)
	t.Log("ratio:", ratio, " err:", err)
	em, err := data.MidValue(_ctx, "EPS", 10, FinaReportTypeYear)
	require.Nil(t, err)
	rm, err := data.MidValue(_ctx, "ROE", 0, FinaReportTypeYear)
	require.Nil(t, err)
	t.Log("eps mid:", em, " roe mid:", rm)
}

func TestQueryFinaPublishDateList(t *testing.T) {
	date, err := _em.QueryFinaPublishDateList(_ctx, "000026")
	require.Nil(t, err)
	t.Log("pubdate:", date)
}
