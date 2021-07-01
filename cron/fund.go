// Package cron 定时任务
package cron

import (
	"context"
	"io/ioutil"
	"sync"
	"time"

	"github.com/avast/retry-go"
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
	reqChan := make(chan string, 100)
	var wg sync.WaitGroup
	var mu sync.Mutex

	fundlist := models.FundList{}
	for _, efund := range efundlist {
		wg.Add(1)
		go func() { reqChan <- efund.Fcode }()
		go func() {
			defer func() {
				wg.Done()
			}()
			code := <-reqChan
			err := retry.Do(
				func() error {
					f, err := datacenter.EastMoney.QueryFundInfo(ctx, code)
					if err != nil {
						return err
					}
					fund, err := models.NewFund(ctx, f)
					if err != nil {
						return err
					}
					mu.Lock()
					fundlist = append(fundlist, fund)
					mu.Unlock()
					return nil
				},
			)
			if err != nil {
				logging.Errorf(ctx, "QueryAllFundList QueryFundInfo err:%v", err)
				promSyncError.WithLabelValues("SyncFundAllList").Inc()
			}
		}()
	}
	wg.Wait()
	logging.Infof(ctx, "SyncFundAllList request end. latency:%+v", time.Now().Sub(start))

	// 更新 services 变量
	services.FundAllList = fundlist

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
