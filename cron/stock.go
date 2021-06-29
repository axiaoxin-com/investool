// Package cron 定时任务
package cron

import (
	"context"
	"io/ioutil"

	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/x-stock/datacenter"
	"github.com/axiaoxin-com/x-stock/services"
	jsoniter "github.com/json-iterator/go"
)

// SyncIndustryList 同步行业列表
func SyncIndustryList() {
	ctx := context.Background()
	indlist, err := datacenter.EastMoney.QueryIndustryList(ctx)
	if err != nil {
		logging.Errorf(ctx, "SyncIndustryList QueryIndustryList error:", err)
		promSyncError.WithLabelValues("SyncIndustryList").Inc()
		return
	}
	services.StockIndustryList = indlist

	// 更新文件
	b, err := jsoniter.Marshal(indlist)
	if err != nil {
		logging.Errorf(ctx, "SyncIndustryList json marshal error:", err)
		promSyncError.WithLabelValues("SyncIndustryList").Inc()
		return
	}
	if err := ioutil.WriteFile(services.IndustryListFilename, b, 0666); err != nil {
		logging.Errorf(ctx, "SyncIndustryList WriteFile error:", err)
		promSyncError.WithLabelValues("SyncIndustryList").Inc()
		return
	}
}
