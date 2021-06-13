// Package services 加载或初始化外部依赖服务
package services

import (
	"context"

	"github.com/axiaoxin-com/x-stock/datacenter"
)

// StockIndustryList 东方财富股票行业列表
var StockIndustryList []string

// Init 相关依赖服务的初始化或加载操作
func Init() error {
	// 初始化行业列表
	indlist, err := datacenter.EastMoney.QueryIndustryList(context.Background())
	if err != nil {
		return err
	}
	StockIndustryList = indlist
	return nil
}
