package core

import (
	"testing"

	"github.com/axiaoxin-com/x-stock/model"
)

func TestCheckFundamentals(t *testing.T) {
	stock := model.Stock{}
	c := NewChecker(_ctx, stock)
	result := c.CheckFundamentals(_ctx)
	t.Log(result)
}
