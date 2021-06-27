package eastmoney

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func _TestQueryAllFundNetList(t *testing.T) {
	data, err := _em.QueryAllFundList(_ctx, FundTypeALL)
	require.Nil(t, err)
	require.NotEmpty(t, data)
	t.Log("data len:", len(data))
}
