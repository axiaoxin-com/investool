package reportor

import (
	"context"
	"testing"
)

var (
	_ctx = context.TODO()
)

func TestGenSelectedStocksReport(t *testing.T) {
	GenSelectedStocksReport(_ctx, "./test.xlsx")
}
