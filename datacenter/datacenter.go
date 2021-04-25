// Package datacenter 数据来源
package datacenter

import (
	"github.com/axiaoxin-com/x-stock/datacenter/eastmoney"
	"github.com/axiaoxin-com/x-stock/datacenter/eniu"
	"github.com/axiaoxin-com/x-stock/datacenter/qq"
)

var (
	// EastMoney 东方财富
	EastMoney eastmoney.EastMoney
	// Eniu 亿牛网
	Eniu eniu.Eniu
	// QQ 腾讯证券
	QQ qq.QQ
)

func init() {
	EastMoney = eastmoney.NewEastMoney()
	Eniu = eniu.NewEniu()
	QQ = qq.NewQQ()
}
