// 关键词搜索

package sina

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"github.com/piex/transcode"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// SearchResult 搜索结果
type SearchResult struct {
	// 数字代码
	SecurityCode string
	// 带后缀的代码
	Secucode string
	// 股票名称
	Name string
}

// KeywordSearch 关键词搜索， 股票、代码、拼音
func (q Sina) KeywordSearch(ctx context.Context, kw string) (results []SearchResult, err error) {
	apiurl := fmt.Sprintf("https://suggest3.sinajs.cn/suggest/key=%s", kw)
	logging.Debug(ctx, "Sina KeywordSearch "+apiurl+" begin")
	beginTime := time.Now()
	resp, err := goutils.HTTPGETRaw(ctx, q.HTTPClient, apiurl)
	utf8resp := transcode.FromString(string(resp)).Decode("GBK").ToString()
	latency := time.Now().Sub(beginTime).Milliseconds()
	logging.Debug(ctx, "Sina KeywordSearch "+apiurl+" end", zap.Int64("latency(ms)", latency), zap.Any("resp", utf8resp))
	if err != nil {
		return nil, err
	}
	ds := strings.Split(utf8resp, "=")
	if len(ds) != 2 {
		return nil, errors.New("search resp invalid:" + utf8resp)
	}
	data := strings.Trim(ds[1], `"`)
	for _, line := range strings.Split(data, ";") {
		lineitems := strings.Split(line, ",")
		if len(lineitems) != 9 {
			continue
		}

		secucode := lineitems[3][2:] + "." + lineitems[3][:2]
		result := SearchResult{
			SecurityCode: lineitems[2],
			Secucode:     secucode,
			Name:         lineitems[6],
		}
		results = append(results, result)
	}
	return
}
