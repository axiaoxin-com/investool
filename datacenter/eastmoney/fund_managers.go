// 天天基金获取基金经理列表(web接口)

package eastmoney

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"github.com/corpix/uarand"
	"go.uber.org/zap"
)

// FundManagerInfo 基金经理信息
type FundManagerInfo struct {
	// ID
	ID string
	// 姓名
	Name string
	// 基金公司id
	Jjgsid string
	// 基金公司名称
	Jjgs string
	// 管理规模（亿元）
	Glgm float64
	// 现任基金代码列表
	FundCodes []string
	// 现任基金名称列表
	FundNames []string
	// 累计从业时间(天)
	TotalDays int
	// 代表基金代码
	DbjjCode string
	// 代表基金名称
	DbjjName string
	// 代表基金收益率(任职期间最佳基金回报)
	Dbjjsyl float64
}

// FundMangers 查询基金经理列表（web接口）
// ft（基金类型） all:全部 gp:股票型 hh:混合型 zq:债券型 sy:收益型
// sc（排序字段）abbname:经理名 jjgs:基金公司 totaldays:从业时间 netnav:基金规模 penavgrowth:现任基金最佳回报
// st（排序类型）asc desc
func (e EastMoney) FundMangers(ctx context.Context, ft, sc, st string) ([]*FundManagerInfo, error) {
	beginTime := time.Now()
	header := map[string]string{
		"user-agent": uarand.GetRandom(),
	}
	result := []*FundManagerInfo{}
	index := 1
	for {
		apiurl := fmt.Sprintf(
			"http://fund.eastmoney.com/Data/FundDataPortfolio_Interface.aspx?dt=14&mc=returnjson&ft=%s&pn=20&pi=%d&sc=%s&st=%s",
			ft,
			index,
			sc,
			st,
		)
		logging.Debug(ctx, "EastMoney FundMangers "+apiurl+" begin", zap.Int("index", index))
		resp, err := goutils.HTTPGETRaw(ctx, e.HTTPClient, apiurl, header)
		strresp := string(resp)
		latency := time.Now().Sub(beginTime).Milliseconds()
		logging.Debug(ctx, "EastMoney FundMangers "+apiurl+" end",
			zap.Int64("latency(ms)", latency),
			// zap.Any("resp", strresp),
		)
		if err != nil {
			return nil, err
		}
		reg, err := regexp.Compile(`\[(".+?")\]`)
		if err != nil {
			logging.Error(ctx, "regexp error:"+err.Error())
			return nil, err
		}
		matched := reg.FindAllStringSubmatch(strresp, -1)
		if len(matched) == 0 {
			break
		}

		for _, m := range matched {
			// "30293769","武建刚","80000252","天治基金","006877,006878","天治量化核心精选混合A,天治量化核心精选混合C","306","-18.86%","006877","天治量化核心精选混合A","0.47亿元","-18.86%"
			field, _ := regexp.Compile(`"(.*?)"`)
			fields := field.FindAllStringSubmatch(m[1], -1)
			if len(fields) != 12 {
				logging.Warnf(ctx, "invalid fields len:%v %v", len(fields), m[1])
				continue
			}
			glgm := 0.0
			if fields[10][1] != "" && fields[10][1] != "--" {
				glgmNum := strings.TrimSuffix(fields[10][1], "亿元")
				glgm, err = strconv.ParseFloat(glgmNum, 64)
				if err != nil {
					logging.Warnf(ctx, "parse glgm:%v to float64 error:%v", glgmNum, err)
				}
			}
			totaldays := 0
			if fields[6][1] != "" && fields[6][1] != "--" {
				totaldays, err = strconv.Atoi(fields[6][1])
				if err != nil {
					logging.Warnf(ctx, "parse totaldays:%v to int error:%v", fields[6], err)
				}
			}
			dbjjsyl := 0.0
			if fields[11][1] != "" && fields[11][1] != "--" {
				dbjjsylNum := strings.TrimSuffix(fields[11][1], "%")
				dbjjsyl, err = strconv.ParseFloat(dbjjsylNum, 64)
				if err != nil {
					logging.Warnf(ctx, "parse dbjjsyl:%v to float64 error:%v", dbjjsylNum, err)
				}
			}
			result = append(result, &FundManagerInfo{
				ID:        fields[0][1],
				Name:      fields[1][1],
				Jjgsid:    fields[2][1],
				Jjgs:      fields[3][1],
				Glgm:      glgm,
				FundCodes: strings.Split(fields[4][1], ","),
				FundNames: strings.Split(fields[5][1], ","),
				TotalDays: totaldays,
				DbjjCode:  fields[8][1],
				DbjjName:  fields[9][1],
				Dbjjsyl:   dbjjsyl,
			})
		}
		index++
	}
	latency := time.Now().Sub(beginTime).Milliseconds()
	logging.Debug(
		ctx,
		"EastMoney FundMangers end",
		zap.Int64("latency(ms)", latency),
		zap.Int("resultCount", len(result)),
	)
	return result, nil
}
