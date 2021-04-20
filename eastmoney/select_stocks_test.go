package eastmoney

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSelectStocks(t *testing.T) {
	data, err := _em.SelectStocks(_ctx)
	require.Nil(t, err)
	require.NotEmpty(t, data)
}

func TestSelectStocksWithFilter(t *testing.T) {
	filter := DefaultFilter
	filter.Industry = "白色家电"
	filter.ListingOver5Y = true
	data, err := _em.SelectStocksWithFilter(_ctx, filter)
	require.Nil(t, err)
	b, _ := json.Marshal(data)
	t.Log(string(b))
}
