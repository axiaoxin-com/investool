package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAutoFilterStocks(t *testing.T) {
	result, err := AutoFilterStocks(_ctx)
	require.Nil(t, err)
	t.Log(result)
}
