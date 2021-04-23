package eastmoney

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	_em  = NewEastMoney()
	_ctx = context.TODO()
)

func TestPost(t *testing.T) {
	url := "https://datacenter.eastmoney.com/stock/selection/api/data/get/"
	body := map[string]string{
		"source": "SELECT_SECURITIES",
		"client": "APP",
		"type":   "RPTA_APP_INDUSTRY",
		"sty":    "ALL",
	}
	rsp := map[string]interface{}{}
	err := _em.Post(_ctx, url, body, &rsp)
	require.Nil(t, err)
	require.Greater(t, len(rsp), 0)
}
