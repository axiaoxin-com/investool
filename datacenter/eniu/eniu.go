// 亿牛网数据源封装

package eniu

import (
	"net/http"
	"time"
)

// Eniu 亿牛网数据源
type Eniu struct {
	// http 客户端
	HTTPClient *http.Client
}

// NewEniu 创建 Eniu 实例
func NewEniu() Eniu {
	hc := &http.Client{
		Timeout: 5 * time.Second,
	}
	return Eniu{
		HTTPClient: hc,
	}
}
