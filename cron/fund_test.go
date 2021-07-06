package cron

import (
	"testing"

	"github.com/axiaoxin-com/logging"
)

func _TestSyncFund(t *testing.T) {
	logging.SetLevel("warn")
	SyncFund()
}
