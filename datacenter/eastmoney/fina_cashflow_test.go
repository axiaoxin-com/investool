package eastmoney

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueryFinaCashflowData(t *testing.T) {
	data, err := _em.QueryFinaCashflowData(_ctx, "000958.SZ")
	t.Log("err:", err)
	require.Nil(t, err)

	require.NotEmpty(t, data)
}
