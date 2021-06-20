// Package services 加载或初始化外部依赖服务
package services

import (
	"context"
	"io/ioutil"

	"github.com/axiaoxin-com/x-stock/datacenter"
	"github.com/axiaoxin-com/x-stock/datacenter/eastmoney"
	jsoniter "github.com/json-iterator/go"
)

// StockIndustryList 东方财富股票行业列表
var StockIndustryList []string

// FundNetList 天天基金净值列表
var FundNetList eastmoney.FundNetList

// FundNetListFilename 基金净值列表数据文件
var FundNetListFilename = "./fundnetlist.json"

// Init 相关依赖服务的初始化或加载操作
func Init() error {
	// 初始化行业列表
	indlist, err := datacenter.EastMoney.QueryIndustryList(context.Background())
	if err != nil {
		return err
	}
	StockIndustryList = indlist

	// 从json文件加载基金列表

	fundlist, err := ioutil.ReadFile(FundNetListFilename)
	if err != nil {
		return err
	}
	if err := jsoniter.Unmarshal(fundlist, FundNetList); err != nil {
		return err
	}

	return nil
}
