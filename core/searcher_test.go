package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearch(t *testing.T) {
	s := NewSearcher(_ctx)
	k := []string{"招商", "贵州茅台", "000001"}
	results, err := s.Search(_ctx, k)
	require.Nil(t, err)
	require.Len(t, results, 3)
}
