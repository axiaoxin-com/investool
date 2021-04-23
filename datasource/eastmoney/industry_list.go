// 获取选股器中的行业列表数据

package eastmoney

import (
	"context"
	"fmt"
)

// RespIndustryList 接口返回的 json 结构
type RespIndustryList struct {
	Result struct {
		Count int `json:"count"`
		Pages int `json:"pages"`
		Data  []struct {
			Industry    string `json:"INDUSTRY"`
			FirstLetter string `json:"FIRST_LETTER"`
		} `json:"data"`
	} `json:"result"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// IndustryList 获取行业列表
func (e EastMoney) IndustryList(ctx context.Context) ([]string, error) {
	url := "https://datacenter.eastmoney.com/stock/selection/api/data/get/"
	body := map[string]string{
		"source": "SELECT_SECURITIES",
		"client": "APP",
		"type":   "RPTA_APP_INDUSTRY",
		"sty":    "ALL",
	}
	resp := RespIndustryList{}
	if err := e.Post(ctx, url, body, &resp); err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("%#v", resp)
	}
	result := []string{}
	for _, i := range resp.Result.Data {
		result = append(result, i.Industry)
	}
	return result, nil
}
