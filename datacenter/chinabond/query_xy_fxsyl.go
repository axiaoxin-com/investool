package chinabond

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"github.com/corpix/uarand"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// RespQueryFxsyl 债券收益率接口返回结果
// SeriesData: [ [期限年数, 收益率], ... ]
type RespQueryFxsyl struct {
	YcChartDataList []struct {
		YcDefID    string      `json:"ycDefId"`
		YcDefName  string      `json:"ycDefName"`
		YcYWName   interface{} `json:"ycYWName"`
		Worktime   string      `json:"worktime"`
		SeriesData [][]float64 `json:"seriesData"`
		IsPoint    bool        `json:"isPoint"`
		HyCurve    bool        `json:"hyCurve"`
		Point      bool        `json:"point"`
	} `json:"ycChartDataList"`
	ChartDataList interface{} `json:"chartDataList"`
	UpThrow       int         `json:"upThrow"`
	DownThrow     int         `json:"downThrow"`
	UpOffset      int         `json:"upOffset"`
	DownOffset    int         `json:"downOffset"`
}

// QueryFxsyl 查询指定债券在指定日期的收益率
// treeItemID 为QueryTree中对应债券的id
// date为string格式的指定日期：YYYY-mm-dd
func (c ChinaBond) QueryFxsyl(ctx context.Context, treeItemID, date string) ([][]float64, error) {
	apiurl := fmt.Sprintf(
		"https://yield.chinabond.com.cn/cbweb-mn/yc/searchXyFxsyl?xyzSelect=txy&&workTimes=%s&&dxbj=4&&qxll=1,&&yqqxN=N&&yqqxK=K&&ycDefIds=%s,&&locale=zh_CN",
		date,
		treeItemID,
	)
	logging.Debug(ctx, "ChinaBond QueryFxsyl "+apiurl+" begin")
	beginTime := time.Now()
	resp := RespQueryFxsyl{}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiurl, nil)
	if err != nil {
		return nil, errors.Wrap(err, "QueryFxsyl NewRequestWithContext")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", uarand.GetRandom())
	err = goutils.HTTPPOST(ctx, c.HTTPClient, req, &resp)
	latency := time.Now().Sub(beginTime).Milliseconds()
	logging.Debug(
		ctx,
		"ChinaBond QueryFxsyl "+apiurl+" end",
		zap.Int64("latency(ms)", latency),
		zap.Any("resp", resp),
	)
	if err != nil {
		return nil, err
	}
	if len(resp.YcChartDataList) == 0 {
		return nil, nil
	}

	return resp.YcChartDataList[0].SeriesData, nil
}
