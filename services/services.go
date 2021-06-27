// Package services 加载或初始化外部依赖服务
package services

import (
	"context"
	"io/ioutil"

	"github.com/axiaoxin-com/x-stock/datacenter"
	"github.com/axiaoxin-com/x-stock/models"
	jsoniter "github.com/json-iterator/go"
)

var (
	// StockIndustryList 东方财富股票行业列表
	StockIndustryList []string
	// FundAllList 全量基金列表
	FundAllList models.FundList
	// Fund4433List 满足4433法则的基金列表
	Fund4433List models.FundList
	// FundAllListFilename 基金列表数据文件
	FundAllListFilename = "./fund_all_list.json"
	// Fund4433ListFilename 4433基金列表数据文件
	Fund4433ListFilename = "./fund_4433_list.json"
	// IndustryListFilename 行业列表数据文件
	IndustryListFilename = "./industry_list.json"
)

// Init 相关依赖服务的初始化或加载操作
func Init() error {
	if err := InitIndustryList(); err != nil {
		return err
	}
	if err := InitFundAllList(); err != nil {
		return err
	}
	if err := InitFund4433List(); err != nil {
		return err
	}
	return nil
}

// InitIndustryList 初始化行业列表
func InitIndustryList() error {
	indlist, err := datacenter.EastMoney.QueryIndustryList(context.Background())
	if err != nil {
		return err
	}
	StockIndustryList = indlist
	return nil
}

// InitFundAllList 从json文件加载基金列表
func InitFundAllList() error {
	fundlist, err := ioutil.ReadFile(FundAllListFilename)
	if err != nil {
		return err
	}
	return jsoniter.Unmarshal(fundlist, &FundAllList)
}

// InitFund4433List 从json文件加载基金列表
func InitFund4433List() error {
	fundlist, err := ioutil.ReadFile(Fund4433ListFilename)
	if err != nil {
		return err
	}
	return jsoniter.Unmarshal(fundlist, &Fund4433List)
}
