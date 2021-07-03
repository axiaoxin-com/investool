package cron

import (
	"testing"

	"github.com/axiaoxin-com/logging"
)

func TestSyncFundAllList(t *testing.T) {
	logging.SetLevel("warn")
	SyncFund()
}
