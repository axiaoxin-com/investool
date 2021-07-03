package cron

import (
	"testing"

	"github.com/axiaoxin-com/logging"
)

func _TestSyncFundAllList(t *testing.T) {
	logging.SetLevel("warn")
	SyncFund()
}
