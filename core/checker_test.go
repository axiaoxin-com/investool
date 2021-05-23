package core

import (
	"testing"

	"github.com/axiaoxin-com/x-stock/model"
)

func TestCheckFundamentals(t *testing.T) {
	stock := model.Stock{}
	c := NewChecker(_ctx, DefaultCheckerOptions)
	result, ok := c.CheckFundamentals(_ctx, stock)
	t.Log(ok, result)
}
