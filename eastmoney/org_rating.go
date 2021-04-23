// 获取机构评级统计

package eastmoney

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"go.uber.org/zap"
)

// RespOrgRating 统计评级接口返回结构
type RespOrgRating struct {
	Version string `json:"version"`
	Result  struct {
		Pages int `json:"pages"`
		Data  []struct {
			Secucode         string      `json:"SECUCODE"`
			SecurityTypeCode string      `json:"SECURITY_TYPE_CODE"`
			TradeMarketCode  string      `json:"TRADE_MARKET_CODE"`
			DateTypeCode     int         `json:"DATE_TYPE_CODE"`
			DateType         string      `json:"DATE_TYPE"`
			CompreRatingNum  float64     `json:"COMPRE_RATING_NUM"`
			CompreRating     string      `json:"COMPRE_RATING"`
			RatingOrgNum     int         `json:"RATING_ORG_NUM"`
			RatingBuyNum     int         `json:"RATING_BUY_NUM"`
			RatingAddNum     int         `json:"RATING_ADD_NUM"`
			RatingNeutralNum interface{} `json:"RATING_NEUTRAL_NUM"`
			RatingReduceNum  interface{} `json:"RATING_REDUCE_NUM"`
			RatingSaleNum    interface{} `json:"RATING_SALE_NUM"`
		} `json:"data"`
		Count int `json:"count"`
	} `json:"result"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// OrgRating 机构评级统计
type OrgRating struct {
	// 时间段
	DateType string `json:"date_type"`
	// 综合评级
	CompreRating string `json:"compre_rating"`
	// 总家数
	RatingOrgNum int `json:"rating_org_num"`
	// 买入
	RatingBuyNum int `json:"rating_buy_num"`
	// 增持
	RatingAddNum int `json:"rating_add_num"`
	// 中性
	RatingNeutralNum interface{} `json:"rating_neutral_num"`
	// 减持
	RatingReduceNum interface{} `json:"rating_reduce_num"`
	// 卖出
	RatingSaleNum interface{} `json:"rating_sale_num"`
}

// QueryOrgRating 获取评级统计
func (e EastMoney) QueryOrgRating(ctx context.Context, secuCode string) ([]OrgRating, error) {
	apiurl := "https://datacenter.eastmoney.com/securities/api/data/get"
	params := map[string]string{
		"source": "SECURITIES",
		"client": "APP",
		"type":   "RPT_RES_ORGRATING",
		"sty":    "SECUCODE,SECURITY_TYPE_CODE,TRADE_MARKET_CODE,DATE_TYPE_CODE,DATE_TYPE,COMPRE_RATING_NUM,COMPRE_RATING,RATING_ORG_NUM,RATING_BUY_NUM,RATING_ADD_NUM,RATING_NEUTRAL_NUM,RATING_REDUCE_NUM,RATING_SALE_NUM",
		"filter": fmt.Sprintf(`(SECUCODE="%s")`, strings.ToUpper(secuCode)),
		"sr":     "1",
		"st":     "DATE_TYPE_CODE",
	}
	logging.Debug(ctx, "EastMoney QueryOrgRating "+apiurl+" begin", zap.Any("params", params))
	beginTime := time.Now()
	apiurl, err := goutils.NewHTTPGetURLWithQueryString(ctx, apiurl, params)
	if err != nil {
		return nil, err
	}
	resp := RespOrgRating{}
	if err := goutils.HTTPGET(ctx, e.HTTPClient, apiurl, &resp); err != nil {
		return nil, err
	}
	latency := time.Now().Sub(beginTime).Milliseconds()
	logging.Debug(ctx, "EastMoney QueryOrgRating "+apiurl+" end", zap.Int64("latency(ms)", latency), zap.Any("resp", resp))
	if resp.Code != 0 {
		return nil, fmt.Errorf("%#v", resp)
	}
	result := []OrgRating{}
	for _, i := range resp.Result.Data {
		or := OrgRating{
			DateType:         i.DateType,
			CompreRating:     i.CompreRating,
			RatingOrgNum:     i.RatingOrgNum,
			RatingBuyNum:     i.RatingBuyNum,
			RatingAddNum:     i.RatingAddNum,
			RatingNeutralNum: i.RatingNeutralNum,
			RatingReduceNum:  i.RatingReduceNum,
			RatingSaleNum:    i.RatingSaleNum,
		}
		result = append(result, or)
	}
	return result, nil
}
