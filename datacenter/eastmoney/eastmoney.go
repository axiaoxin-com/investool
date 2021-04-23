// 东方财富数据源封装

package eastmoney

import (
	"net/http"
	"time"
)

// EastMoney 东方财富数据源
type EastMoney struct {
	// http 客户端
	HTTPClient *http.Client
}

// NewEastMoney 创建 EastMoney 实例
func NewEastMoney() EastMoney {
	hc := &http.Client{
		Timeout: 60 * time.Second,
	}
	return EastMoney{
		HTTPClient: hc,
	}
}
