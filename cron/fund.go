// Package cron 定时任务
package cron

import (
	"context"
	"io/ioutil"
	"time"

	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/x-stock/core"
	"github.com/axiaoxin-com/x-stock/datacenter"
	"github.com/axiaoxin-com/x-stock/datacenter/eastmoney"
	"github.com/axiaoxin-com/x-stock/models"
	"github.com/axiaoxin-com/x-stock/services"
	jsoniter "github.com/json-iterator/go"
)

// SyncFund 同步基金数据
func SyncFund() {
	ctx := context.Background()
	logging.Infof(ctx, "SyncFund request start...")

	// 获取全量列表
	efundlist, err := datacenter.EastMoney.QueryAllFundList(ctx, eastmoney.FundTypeALL)
	if err != nil {
		logging.Error(ctx, "SyncFund QueryAllFundList error:"+err.Error())
		promSyncError.WithLabelValues("SyncFund").Inc()
		return
	}

	fundCodes := []string{}
	for _, efund := range efundlist {
		fundCodes = append(fundCodes, efund.Fcode)
	}
	s := core.NewSearcher(ctx)
	data, err := s.SearchFunds(ctx, fundCodes)
	fundlist := models.FundList{}
	typeMap := map[string]struct{}{}
	for _, fund := range data {
		fundlist = append(fundlist, fund)
		typeMap[fund.Type] = struct{}{}
	}

	// 更新 services 变量
	services.FundAllList = fundlist
	fundtypes := []string{}
	for k := range typeMap {
		fundtypes = append(fundtypes, k)
	}
	services.FundTypeList = fundtypes

	// 更新文件
	b, err := jsoniter.Marshal(fundlist)
	if err != nil {
		logging.Errorf(ctx, "SyncFund json marshal fundlist error:", err)
		promSyncError.WithLabelValues("SyncFund").Inc()
	}
	if err := ioutil.WriteFile(services.FundAllListFilename, b, 0666); err != nil {
		logging.Errorf(ctx, "SyncFund WriteFile fundlist error:", err)
		promSyncError.WithLabelValues("SyncFund").Inc()
	}
	b, err = jsoniter.Marshal(services.FundTypeList)
	if err != nil {
		logging.Errorf(ctx, "SyncFund json marshal fundtypelist error:", err)
		promSyncError.WithLabelValues("SyncFund").Inc()
	}
	if err := ioutil.WriteFile(services.FundTypeListFilename, b, 0666); err != nil {
		logging.Errorf(ctx, "SyncFund WriteFile fundtypelist error:", err)
		promSyncError.WithLabelValues("SyncFund").Inc()
	}

	// 更新4433列表
	Update4433()

	// 更新同步时间
	services.SyncFundTime = time.Now()
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
	fundlist.Sort(models.FundSortTypeWeek)
	services.Fund4433List = fundlist

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
