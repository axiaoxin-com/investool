package core

import (
	"testing"

	"github.com/axiaoxin-com/x-stock/datacenter/eastmoney"
	"github.com/stretchr/testify/require"
)

func TestAutoFilterStocks(t *testing.T) {
	checker := NewChecker(_ctx, DefaultCheckerOptions)
	s := NewSelector(_ctx, eastmoney.DefaultFilter, &checker)
	result, err := s.AutoFilterStocks(_ctx)
	require.Nil(t, err)
	t.Log(result)
}
