// Package cron 定时任务
package cron

import (
	"context"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/investool/datacenter"
	"github.com/axiaoxin-com/investool/services"
)

// SyncBond 同步债券
func SyncBond() {
	if !goutils.IsTradingDay() {
		return
	}
	ctx := context.Background()
	services.AAACompanyBondSyl = datacenter.ChinaBond.QueryAAACompanyBondSyl(ctx)
}
