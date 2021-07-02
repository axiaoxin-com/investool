package models

import (
	"context"
	"testing"

	"github.com/axiaoxin-com/x-stock/datacenter/eastmoney"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
)

func TestNewFund(t *testing.T) {
	ctx := context.TODO()
	efund, err := eastmoney.NewEastMoney().QueryFundInfo(ctx, "260104")
	require.Nil(t, err)
	fund := NewFund(ctx, efund)
	b, err := jsoniter.Marshal(fund)
	require.Nil(t, err)
	t.Log(string(b))
}
