package core

import (
	"testing"

	"github.com/axiaoxin-com/x-stock/services"
)

func TestFunderFilter(t *testing.T) {
	services.FundAllListFilename = "../fund_all_list.json"
	services.InitFundAllList()
	funder := NewFunder()
	p := ParamFunderFilter{
		Types:                []string{},
		MinScale:             2.0,
		MaxScale:             50.0,
		MinManagerYears:      5.0,
		Year1RankRatio:       25.0,
		ThisYear235RankRatio: 25.0,
		Month6RankRatio:      30.0,
		Month3RankRatio:      30.0,
		Max135AvgStddev:      0.0,
		Min135AvgSharp:       0.0,
		Max135AvgRetr:        0.0,
	}
	results := funder.Filter(_ctx, p)
	t.Log(len(results))
}
