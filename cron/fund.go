// Package cron 定时任务
package cron

import (
	"context"
	"io/ioutil"
	"math"
	"sync"
	"time"

	"github.com/avast/retry-go"
	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/x-stock/datacenter"
	"github.com/axiaoxin-com/x-stock/datacenter/eastmoney"
	"github.com/axiaoxin-com/x-stock/models"
	"github.com/axiaoxin-com/x-stock/services"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
)

// SyncFund 同步基金数据
func SyncFund() {
	ctx := context.Background()
	start := time.Now()
	logging.Infof(ctx, "SyncFund request start...")

	// 获取全量列表
	efundlist, err := datacenter.EastMoney.QueryAllFundList(ctx, eastmoney.FundTypeALL)
	if err != nil {
		logging.Errorf(ctx, "SyncFund QueryAllFundList error:", err)
		promSyncError.WithLabelValues("SyncFund").Inc()
		return
	}

	// 遍历获取基金详情
	chanSize := viper.GetFloat64("app.chan_size")
	if chanSize == 0 {
		chanSize = 500
	}
	workerCount := int(math.Min(float64(len(efundlist)), chanSize))
	reqChan := make(chan string, workerCount)
	typeMap := map[string]struct{}{}
	var wg sync.WaitGroup
	var mu sync.Mutex

	fundlist := models.FundList{}
	for _, efund := range efundlist {
		wg.Add(1)
		reqChan <- efund.Fcode
		go func() {
			defer func() {
				wg.Done()
			}()
			// 低配机器 oom fix
			time.Sleep(time.Millisecond * 100)

			code := <-reqChan
			fundresp := &eastmoney.RespFundInfo{}
			err := retry.Do(
				func() error {
					var err error
					fundresp, err = datacenter.EastMoney.QueryFundInfo(ctx, code)
					return err

				},
				retry.OnRetry(func(n uint, err error) {
					logging.Errorf(ctx, "retry#%d: code:%v %v\n", n, code, err)
				}),
			)
			if err != nil {
				logging.Errorf(ctx, "QueryAllFundList QueryFundInfo code:%v err:%v", code, err)
				promSyncError.WithLabelValues("SyncFund").Inc()
				return
			}
			fund := models.NewFund(ctx, fundresp)
			mu.Lock()
			fundlist = append(fundlist, fund)
			typeMap[fund.Type] = struct{}{}
			mu.Unlock()
		}()
	}
	wg.Wait()
	logging.Infof(ctx, "SyncFund request end. latency:%+v", time.Now().Sub(start))

	// 更新 services 变量
	services.FundAllList = fundlist
	services.FundTypeList = []string{}
	for k := range typeMap {
		services.FundTypeList = append(services.FundTypeList, k)
	}

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
