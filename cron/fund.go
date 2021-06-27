// Package cron 定时任务
package cron

import (
	"context"
	"io/ioutil"
	"time"

	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/x-stock/datacenter"
	"github.com/axiaoxin-com/x-stock/datacenter/eastmoney"
	"github.com/axiaoxin-com/x-stock/models"
	"github.com/axiaoxin-com/x-stock/services"
	jsoniter "github.com/json-iterator/go"
)

// SyncFundAllList 同步基金列表
func SyncFundAllList() {
	ctx := context.Background()
	start := time.Now()
	logging.Infof(ctx, "SyncFundAllList request start...")

	// 获取全量列表
	efundlist, err := datacenter.EastMoney.QueryAllFundList(ctx, eastmoney.FundTypeALL)
	if err != nil {
		logging.Errorf(ctx, "SyncFundAllList QueryAllFundList error:", err)
		promSyncError.WithLabelValues("SyncFundAllList").Inc()
		return
	}

	// 遍历获取基金详情
	fundlist := models.FundList{}
	for _, efund := range efundlist {
		f, err := datacenter.EastMoney.QueryFundInfo(ctx, efund.Fcode)
		if err != nil {
			logging.Errorf(ctx, "SyncFundAllList QueryFundInfo error:%v", err)
			promSyncError.WithLabelValues("SyncFundAllList").Inc()
			continue
		}
		fund, err := models.NewFund(ctx, f)
		if err != nil {
			logging.Errorf(ctx, "SyncFundAllList NewFund error:%v", err)
			promSyncError.WithLabelValues("SyncFundAllList").Inc()
			continue
		}
		fundlist = append(fundlist, fund)
	}
	logging.Infof(ctx, "SyncFundAllList request end. latency:%+v", time.Now().Sub(start))

	// 更新 services 变量
	rwMutex.RLock()
	services.FundAllList = fundlist
	rwMutex.RUnlock()

	// 更新文件
	b, err := jsoniter.Marshal(fundlist)
	if err != nil {
		logging.Errorf(ctx, "SyncFundAllList json marshal error:", err)
		promSyncError.WithLabelValues("SyncFundAllList").Inc()
		return
	}
	if err := ioutil.WriteFile(services.FundAllListFilename, b, 0666); err != nil {
		logging.Errorf(ctx, "SyncFundAllList WriteFile error:", err)
		promSyncError.WithLabelValues("SyncFundAllList").Inc()
		return
	}
}

// Update4433 更新4433检测结果
func Update4433() {
	ctx := context.Background()
	fundlist := models.FundList{}
	for _, fund := range services.FundAllList {
		if fund.Is4433(ctx) {
			fundlist = append(fundlist, fund)
		}
	}
	// 更新 services 变量
	rwMutex.RLock()
	services.Fund4433List = fundlist
	rwMutex.RUnlock()

	// 更新文件
	b, err := jsoniter.Marshal(fundlist)
	if err != nil {
		logging.Errorf(ctx, "Update4433 json marshal error:", err)
		promSyncError.WithLabelValues("Update4433").Inc()
		return
	}
	if err := ioutil.WriteFile(services.Fund4433ListFilename, b, 0666); err != nil {
		logging.Errorf(ctx, "Update4433 WriteFile error:", err)
		promSyncError.WithLabelValues("Update4433").Inc()
		return
	}
}
