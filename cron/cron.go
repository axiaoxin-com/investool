// Package cron 定时任务
package cron

import (
	"context"
	"io/ioutil"
	"sync"
	"time"

	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/x-stock/datacenter"
	"github.com/axiaoxin-com/x-stock/datacenter/eastmoney"
	"github.com/axiaoxin-com/x-stock/services"
	"github.com/go-co-op/gocron"
	jsoniter "github.com/json-iterator/go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	rwMutex = new(sync.RWMutex)

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
	sched.Cron("0 3 * * 6").Do(SyncFundNetList)

	// 每月1号凌晨4点同步东方财富行业列表
	sched.Cron("0 4 1 * *").Do(SyncIndustryList)
}

// SyncFundNetList 同步基金净值列表
func SyncFundNetList() {
	ctx := context.Background()
	fundlist, err := datacenter.EastMoney.QueryAllFundNetList(ctx, eastmoney.FundTypeALL)
	if err != nil {
		logging.Errorf(ctx, "SyncFundNetList QueryAllFundNetList error:", err)
		promSyncError.WithLabelValues("SyncFundNetList").Inc()
		return
	}

	// 更新 services 变量
	rwMutex.RLock()
	services.FundNetList = fundlist
	rwMutex.RUnlock()

	// 更新文件
	b, err := jsoniter.Marshal(fundlist)
	if err != nil {
		logging.Errorf(ctx, "SyncFundNetList json marshal error:", err)
		promSyncError.WithLabelValues("SyncFundNetList").Inc()
		return
	}
	if err := ioutil.WriteFile(services.FundNetListFilename, b, 0666); err != nil {
		logging.Errorf(ctx, "SyncFundNetList WriteFile error:", err)
		promSyncError.WithLabelValues("SyncFundNetList").Inc()
		return
	}
}

// SyncIndustryList 同步行业列表
func SyncIndustryList() {
	ctx := context.Background()
	indlist, err := datacenter.EastMoney.QueryIndustryList(ctx)
	if err != nil {
		logging.Errorf(ctx, "SyncIndustryList QueryIndustryList error:", err)
		promSyncError.WithLabelValues("SyncIndustryList").Inc()
		return
	}
	rwMutex.RLock()
	services.StockIndustryList = indlist
	rwMutex.RUnlock()
}
