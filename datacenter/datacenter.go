// Package datacenter 数据来源
package datacenter

import (
	"github.com/axiaoxin-com/x-stock/datacenter/eastmoney"
	"github.com/axiaoxin-com/x-stock/datacenter/eniu"
	"github.com/axiaoxin-com/x-stock/datacenter/sina"
	"github.com/axiaoxin-com/x-stock/datacenter/zszx"
)

var (
	// EastMoney 东方财富
	EastMoney eastmoney.EastMoney
	// Eniu 亿牛网
	Eniu eniu.Eniu
	// Sina 新浪财经
	Sina sina.Sina
	// Zszx 招商证券
	Zszx zszx.Zszx
)

func init() {
	EastMoney = eastmoney.NewEastMoney()
	Eniu = eniu.NewEniu()
	Sina = sina.NewSina()
	Zszx = zszx.NewZszx()
}
