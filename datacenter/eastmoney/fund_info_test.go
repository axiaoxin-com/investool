package eastmoney

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueryFundInfo(t *testing.T) {
	data, err := _em.QueryFundInfo(_ctx, "007135")
	require.Nil(t, err)
	require.NotEmpty(t, data)
	t.Log("data:", data)
}
