// Package cron 定时任务
package cron

import (
	"time"

	"github.com/go-co-op/gocron"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	promSyncLabels = []string{
		"jobname",
	}
	promSyncError = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "cron",
			Name:      "sync_error",
			Help:      "cron sync job error",
		}, promSyncLabels,
	)
)

// RunCronJobs 启动定时任务
func RunCronJobs() {
	timezone, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	sched := gocron.NewScheduler(timezone)
	// 每周六凌晨3点同步基金净值列表
	sched.Cron("0 3 * * 6").Do(SyncFundAllList)
	// 每周六凌晨6点更新4433基金列表
	sched.Cron("0 6 * * 6").Do(Update4433)
	// 每月1号凌晨4点同步东方财富行业列表
	sched.Cron("0 4 1 * *").Do(SyncIndustryList)
}
