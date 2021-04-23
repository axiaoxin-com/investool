package eastmoney

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueryValuationStatus(t *testing.T) {
	status, data, err := _em.QueryValuationStatus(_ctx, "600149.sh")
	require.Nil(t, err)
	require.Len(t, data, 4)
	t.Log("status:", status, " data:", data)
}
