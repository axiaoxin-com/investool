// 天天基金获取基金经理列表

package eastmoney

import (
	"context"
	"fmt"
	"time"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"github.com/corpix/uarand"
	"go.uber.org/zap"
)

// FundManagerBaseInfo 基金经理基本信息
type FundManagerBaseInfo struct {
	// ID
	Mgrid string `json:"MGRID"`
	// 姓名
	Mgrname string `json:"MGRNAME"`
	// 擅长领域 1:偏债类 2:偏股类 3:指数类 4:货币类 5:QDII
	Mftype string `json:"MFTYPE"`
	// 基金公司
	Jjgs string `json:"JJGS"`
	// 基金公司id
	Jjgsid string `json:"JJGSID"`
	// 从业年化收益（%）
	Yieldse string `json:"YIELDSE"`
	// 近1周涨跌幅（%）
	W string `json:"W"`
	// 近1月涨跌幅（%）
	M string `json:"M"`
	// 近3月涨跌幅（%）
	Q string `json:"Q"`
	// 近6月涨跌幅（%）
	Hy string `json:"HY"`
	// 近1年涨跌幅（%）
	Y string `json:"Y"`
	// 管理规模（元）
	Netnav string `json:"NETNAV"`
	// 业绩评分
	Mgold string `json:"MGOLD"`
	// 代表基金代码
	Precode string `json:"PRECODE"`
	// 代表基金名称
	Shortname string `json:"SHORTNAME"`
	// 头像
	Newphotourl string `json:"NEWPHOTOURL"`
	// 性别 0:男 1:女
	Sex string `json:"SEX"`
}

// RespFundMangerBaseList FundMangerBaseList 接口原始返回结构
type RespFundMangerBaseList struct {
	Datas      []*FundManagerBaseInfo `json:"Datas"`
	ErrCode    int                    `json:"ErrCode"`
	ErrMsg     interface{}            `json:"ErrMsg"`
	TotalCount int                    `json:"TotalCount"`
	Expansion  interface{}            `json:"Expansion"`
}

// FundMangerBaseList 查询基金经理列表（app接口）
// mftype "":全部 1:偏债类 2:偏股类 3:指数类 4:货币类 5:QDII
// sortColum W:近1周平均收益 M:近1月平均收益 Q:近3月平均收益 HY:近半年平均收益 Y:近1年平均收益 NETNAV:管理规模 MGOLD:业绩评价 YIELDSE:从业年化回报
func (e EastMoney) FundMangerBaseList(ctx context.Context, mftype string, sortColum string) ([]*FundManagerBaseInfo, error) {
	beginTime := time.Now()
	resp := RespFundMangerBaseList{}
	header := map[string]string{
		"user-agent": uarand.GetRandom(),
	}
	result := []*FundManagerBaseInfo{}
	index := 1
	size := 300
	total := 0
	for {
		apiurl := fmt.Sprintf(
			"https://fundmapi.eastmoney.com/fundmobapi/FundMApi/FundMangerBaseList.ashx?COMPANYCODES=&MFTYPE=%s&Sort=desc&SortColumn=%s&deviceid=fundmanager2016&pageIndex=%d&pageSize=%d&plat=Iphone&product=EFund&version=4.3.0",
			mftype,
			sortColum,
			index,
			size,
		)
		logging.Debug(ctx, "EastMoney FundMangerBaseList "+apiurl+" begin", zap.Int("index", index))
		if err := goutils.HTTPGET(ctx, e.HTTPClient, apiurl, header, &resp); err != nil {
			return nil, err
		}
		total = resp.TotalCount
		if len(resp.Datas) == 0 {
			break
		}
		result = append(result, resp.Datas...)
		index++
	}
	latency := time.Now().Sub(beginTime).Milliseconds()
	logging.Debug(
		ctx,
		"EastMoney FundMangerBaseList end",
		zap.Int64("latency(ms)", latency),
		zap.Int("totalCount", total),
		zap.Int("resultCount", len(result)),
	)
	return result, nil
}
