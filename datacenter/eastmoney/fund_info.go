// 天天基金获取基金详情

package eastmoney

import (
	"context"
	"fmt"
	"time"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"go.uber.org/zap"
)

// QueryFundInfo 查询基金详情
func (e EastMoney) QueryFundInfo(ctx context.Context, fundCode string) (RespFundNetList, error) {
	apiurl := "https://fundmobapi.eastmoney.com/FundMNewApi/FundMNNetNewList"
	params := map[string]string{}
	logging.Debug(ctx, "EastMoney QueryFundInfo "+apiurl+" begin", zap.Any("params", params))
	beginTime := time.Now()
	apiurl, err := goutils.NewHTTPGetURLWithQueryString(ctx, apiurl, params)
	if err != nil {
		return RespFundNetList{}, err
	}
	resp := RespFundNetList{}
	err = goutils.HTTPGET(ctx, e.HTTPClient, apiurl, &resp)
	latency := time.Now().Sub(beginTime).Milliseconds()
	logging.Debug(
		ctx,
		"EastMoney QueryFundInfo "+apiurl+" end",
		zap.Int64("latency(ms)", latency),
		zap.Any("resp", resp),
	)
	if err != nil {
		return resp, err
	}
	if resp.ErrCode != 0 {
		return resp, fmt.Errorf("QueryFundInfo error %v", resp.ErrMsg)
	}
	return resp, nil
}
