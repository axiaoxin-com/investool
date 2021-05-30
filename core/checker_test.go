package core

import (
	"testing"

	"github.com/axiaoxin-com/x-stock/models"
)

func TestCheckFundamentals(t *testing.T) {
	stock := models.Stock{}
	c := NewChecker(_ctx, DefaultCheckerOptions)
	result, ok := c.CheckFundamentals(_ctx, stock)
	t.Log(ok, result)
}
