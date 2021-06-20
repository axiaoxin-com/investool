package eastmoney

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueryAllFundNetList(t *testing.T) {
	data, err := _em.QueryAllFundNetList(_ctx, FundTypeALL)
	require.Nil(t, err)
	require.NotEmpty(t, data)
	t.Log("data len:", len(data))
}
