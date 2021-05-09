package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAutoFilterStocks(t *testing.T) {
	f := NewSelector(_ctx)
	result, err := f.AutoFilterStocks(_ctx, true)
	require.Nil(t, err)
	t.Log(result)
}
