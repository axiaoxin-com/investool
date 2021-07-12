package eastmoney

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func _TestQueryAllFundNetList(t *testing.T) {
	viper.SetDefault("app.chan_size", 500)
	data, err := _em.QueryAllFundList(_ctx, FundTypeALL)
	require.Nil(t, err)
	require.NotEmpty(t, data)
	t.Log("data len:", len(data))
}
