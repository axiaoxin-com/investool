package eastmoney

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSelectStocks(t *testing.T) {
	data, err := _em.SelectStocks(_ctx)
	require.Nil(t, err)
	require.Equal(t, 98, len(data))
	data.SortByROE()
	b, _ := json.Marshal(data)
	t.Log(string(b))
}

func TestSelectStocksWithFilter(t *testing.T) {
	filter := DefaultFilter
	filter.Industry = "玻璃"
	data, err := _em.SelectStocksWithFilter(_ctx, filter)
	require.Nil(t, err)
	require.Equal(t, 1, len(data))
	b, _ := json.Marshal(data)
	t.Log(string(b))
}
