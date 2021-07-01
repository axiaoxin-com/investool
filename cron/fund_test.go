package cron

import (
	"testing"

	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/x-stock/services"
	"github.com/stretchr/testify/require"
)

func _TestSyncFundAllList(t *testing.T) {
	logging.SetLevel("error")
	SyncFundAllList()
}

func _TestUpdate4433(t *testing.T) {
	services.FundAllListFilename = "../fund_all_list.json"
	err := services.InitFundAllList()
	require.Nil(t, err)
	Update4433()
}
