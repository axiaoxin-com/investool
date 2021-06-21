// 天天基金获取基金详情

package eastmoney

import (
	"context"
	"fmt"
	"time"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"go.uber.org/zap"
)

// RespFundInfo QueryFundInfo 接口原始返回结构
type RespFundInfo struct {
	Jjxq struct {
		Datas struct {
			Fcode              string `json:"FCODE"`
			Shortname          string `json:"SHORTNAME"`
			Ftype              string `json:"FTYPE"`
			Feature            string `json:"FEATURE"`
			Bfundtype          string `json:"BFUNDTYPE"`
			Fundtype           string `json:"FUNDTYPE"`
			Rzdf               string `json:"RZDF"`
			Dwjz               string `json:"DWJZ"`
			Ljjz               string `json:"LJJZ"`
			Sgzt               string `json:"SGZT"`
			Shzt               string `json:"SHZT"`
			Sourcerate         string `json:"SOURCERATE"`
			Rate               string `json:"RATE"`
			Minsg              string `json:"MINSG"`
			Maxsg              string `json:"MAXSG"`
			Subscribetime      string `json:"SUBSCRIBETIME"`
			Risklevel          string `json:"RISKLEVEL"`
			Isbuy              string `json:"ISBUY"`
			Bagtype            string `json:"BAGTYPE"`
			Cashbuy            string `json:"CASHBUY"`
			Saletocash         string `json:"SALETOCASH"`
			Stktocash          string `json:"STKTOCASH"`
			Stkexchg           string `json:"STKEXCHG"`
			Fundexchg          string `json:"FUNDEXCHG"`
			Buy                bool   `json:"BUY"`
			Issales            string `json:"ISSALES"`
			Salemark           string `json:"SALEMARK"`
			Mindt              string `json:"MINDT"`
			Dtzt               string `json:"DTZT"`
			Realsgcode         string `json:"REALSGCODE"`
			Qdtcode            string `json:"QDTCODE"`
			Backcode           string `json:"BACKCODE"`
			Estabdate          string `json:"ESTABDATE"`
			Indexcode          string `json:"INDEXCODE"`
			Indexname          string `json:"INDEXNAME"`
			Indextexch         string `json:"INDEXTEXCH"`
			Newindextexch      string `json:"NEWINDEXTEXCH"`
			RlevelSz           string `json:"RLEVEL_SZ"`
			Sharp1             string `json:"SHARP1"`
			Sharp2             string `json:"SHARP2"`
			Sharp3             string `json:"SHARP3"`
			Maxretra1          string `json:"MAXRETRA1"`
			Stddev1            string `json:"STDDEV1"`
			Stddev2            string `json:"STDDEV2"`
			Stddev3            string `json:"STDDEV3"`
			Ssbcfmdata         string `json:"SSBCFMDATA"`
			Ssbcfday           string `json:"SSBCFDAY"`
			Currentdaymark     string `json:"CURRENTDAYMARK"`
			Buymark            string `json:"BUYMARK"`
			Jjgs               string `json:"JJGS"`
			Jjgsid             string `json:"JJGSID"`
			Tsrq               string `json:"TSRQ"`
			Ttypename          string `json:"TTYPENAME"`
			Ttype              string `json:"TTYPE"`
			FundSubjectURL     string `json:"FundSubjectURL"`
			Fbkindexcode       string `json:"FBKINDEXCODE"`
			Fbkindexname       string `json:"FBKINDEXNAME"`
			Fsrq               string `json:"FSRQ"`
			Issbdate           string `json:"ISSBDATE"`
			Rgbegin            string `json:"RGBEGIN"`
			Issedate           string `json:"ISSEDATE"`
			Rgend              string `json:"RGEND"`
			Listtexch          string `json:"LISTTEXCH"`
			Newtexch           string `json:"NEWTEXCH"`
			Islist             string `json:"ISLIST"`
			Islisttrade        string `json:"ISLISTTRADE"`
			Minsbsg            string `json:"MINSBSG"`
			Minsbrg            string `json:"MINSBRG"`
			Endnav             string `json:"ENDNAV"`
			Fegmrq             string `json:"FEGMRQ"`
			Isfnew             string `json:"ISFNEW"`
			Isappoint          string `json:"ISAPPOINT"`
			Minrg              string `json:"MINRG"`
			Cycle              string `json:"CYCLE"`
			Opestart           string `json:"OPESTART"`
			Opeend             string `json:"OPEEND"`
			Opeyield           string `json:"OPEYIELD"`
			Fixincome          string `json:"FIXINCOME"`
			Appointment        string `json:"APPOINTMENT"`
			Appointmenturl     string `json:"APPOINTMENTURL"`
			Isabnormal         string `json:"ISABNORMAL"`
			Yzba               string `json:"YZBA"`
			Fbyzq              string `json:"FBYZQ"`
			Kfsgsh             string `json:"KFSGSH"`
			Linkzsb            string `json:"LINKZSB"`
			Listtexchmark      string `json:"LISTTEXCHMARK"`
			Isharebonus        bool   `json:"ISHAREBONUS"`
			PtdtY              string `json:"PTDT_Y"`
			PtdtTwy            string `json:"PTDT_TWY"`
			PtdtTry            string `json:"PTDT_TRY"`
			PtdtFy             string `json:"PTDT_FY"`
			MbdtY              string `json:"MBDT_Y"`
			MbdtTwy            string `json:"MBDT_TWY"`
			MbdtTry            string `json:"MBDT_TRY"`
			MbdtFy             string `json:"MBDT_FY"`
			YddtY              string `json:"YDDT_Y"`
			YddtTwy            string `json:"YDDT_TWY"`
			YddtTry            string `json:"YDDT_TRY"`
			YddtFy             string `json:"YDDT_FY"`
			DwdtY              string `json:"DWDT_Y"`
			DwdtTwy            string `json:"DWDT_TWY"`
			DwdtTry            string `json:"DWDT_TRY"`
			DwdtFy             string `json:"DWDT_FY"`
			Isyydt             string `json:"ISYYDT"`
			SylZ               string `json:"SYL_Z"`
			Syrq               string `json:"SYRQ"`
			Comethod           string `json:"COMETHOD"`
			Mcoverdate         string `json:"MCOVERDATE"`
			Mcoverdetail       string `json:"MCOVERDETAIL"`
			Comments           string `json:"COMMENTS"`
			Trkerror           string `json:"TRKERROR"`
			Estdiff            string `json:"ESTDIFF"`
			Hrgrt              string `json:"HRGRT"`
			Hsgrt              string `json:"HSGRT"`
			Bench              string `json:"BENCH"`
			Finsales           string `json:"FINSALES"`
			Investmentidear    string `json:"INVESTMENTIDEAR"`
			Investmentidearimg string `json:"INVESTMENTIDEARIMG"`
		} `json:"Datas"`
		ErrCode      int         `json:"ErrCode"`
		Success      bool        `json:"Success"`
		ErrMsg       interface{} `json:"ErrMsg"`
		Message      interface{} `json:"Message"`
		ErrorCode    string      `json:"ErrorCode"`
		ErrorMessage interface{} `json:"ErrorMessage"`
		ErrorMsgLst  interface{} `json:"ErrorMsgLst"`
		TotalCount   int         `json:"TotalCount"`
		Expansion    interface{} `json:"Expansion"`
	} `json:"JJXQ"`
	Jjbq struct {
		Datas []struct {
			Featype string `json:"FEATYPE"`
			Taglist []struct {
				Feavalue string `json:"FEAVALUE"`
				Featype  string `json:"FEATYPE"`
				Feaname  string `json:"FEANAME"`
				Feabrief string `json:"FEABRIEF"`
				Featag   string `json:"FEATAG"`
				Listcode string `json:"LISTCODE"`
			} `json:"TAGLIST"`
		} `json:"Datas"`
		ErrCode      int         `json:"ErrCode"`
		Success      bool        `json:"Success"`
		ErrMsg       interface{} `json:"ErrMsg"`
		Message      interface{} `json:"Message"`
		ErrorCode    string      `json:"ErrorCode"`
		ErrorMessage interface{} `json:"ErrorMessage"`
		ErrorMsgLst  interface{} `json:"ErrorMsgLst"`
		TotalCount   int         `json:"TotalCount"`
		Expansion    interface{} `json:"Expansion"`
	} `json:"JJBQ"`
	Fhts struct {
		Datas struct {
			SameType struct {
				Fcode     interface{} `json:"FCODE"`
				Shortname interface{} `json:"SHORTNAME"`
				Ftype     interface{} `json:"FTYPE"`
				Fundtype  interface{} `json:"FUNDTYPE"`
				Feature   interface{} `json:"FEATURE"`
				Bfundtype interface{} `json:"BFUNDTYPE"`
				Rzdf      interface{} `json:"RZDF"`
				Dwjz      interface{} `json:"DWJZ"`
				Hldwjz    string      `json:"HLDWJZ"`
				Ljjz      interface{} `json:"LJJZ"`
				Ftyi      interface{} `json:"FTYI"`
				Teyi      interface{} `json:"TEYI"`
				Tfyi      interface{} `json:"TFYI"`
				SylZ      interface{} `json:"SYL_Z"`
				SylY      interface{} `json:"SYL_Y"`
				Syl3Y     string      `json:"SYL_3Y"`
				Syl6Y     interface{} `json:"SYL_6Y"`
				Syl1N     interface{} `json:"SYL_1N"`
				Syl2N     interface{} `json:"SYL_2N"`
				Syl3N     interface{} `json:"SYL_3N"`
				Syl5N     interface{} `json:"SYL_5N"`
				SylJn     interface{} `json:"SYL_JN"`
				SylLn     interface{} `json:"SYL_LN"`
			} `json:"SameType"`
			Rele    interface{} `json:"Rele"`
			Subject interface{} `json:"subject"`
		} `json:"Datas"`
		ErrCode      int         `json:"ErrCode"`
		Success      bool        `json:"Success"`
		ErrMsg       interface{} `json:"ErrMsg"`
		Message      interface{} `json:"Message"`
		ErrorCode    string      `json:"ErrorCode"`
		ErrorMessage interface{} `json:"ErrorMessage"`
		ErrorMsgLst  interface{} `json:"ErrorMsgLst"`
		TotalCount   int         `json:"TotalCount"`
		Expansion    interface{} `json:"Expansion"`
	} `json:"FHTS"`
	Jjtx struct {
		Datas struct {
			Trademarklist []interface{} `json:"TRADEMARKLIST"`
			Warnlist      []interface{} `json:"WARNLIST"`
			Sgztmark      interface{}   `json:"SGZTMARK"`
		} `json:"Datas"`
		ErrCode      int         `json:"ErrCode"`
		Success      bool        `json:"Success"`
		ErrMsg       interface{} `json:"ErrMsg"`
		Message      interface{} `json:"Message"`
		ErrorCode    string      `json:"ErrorCode"`
		ErrorMessage interface{} `json:"ErrorMessage"`
		ErrorMsgLst  interface{} `json:"ErrorMsgLst"`
		TotalCount   int         `json:"TotalCount"`
		Expansion    interface{} `json:"Expansion"`
	} `json:"JJTX"`
	Jdzf struct {
		Datas []struct {
			Title string `json:"title"`
			Syl   string `json:"syl"`
			Avg   string `json:"avg"`
			Hs300 string `json:"hs300"`
			Rank  string `json:"rank"`
			Sc    string `json:"sc"`
			Diff  string `json:"diff"`
		} `json:"Datas"`
		ErrCode      int         `json:"ErrCode"`
		Success      bool        `json:"Success"`
		ErrMsg       interface{} `json:"ErrMsg"`
		Message      interface{} `json:"Message"`
		ErrorCode    string      `json:"ErrorCode"`
		ErrorMessage interface{} `json:"ErrorMessage"`
		ErrorMsgLst  interface{} `json:"ErrorMsgLst"`
		TotalCount   int         `json:"TotalCount"`
		Expansion    struct {
			Estabdate  string `json:"ESTABDATE"`
			Time       string `json:"TIME"`
			Isupdating bool   `json:"ISUPDATING"`
		} `json:"Expansion"`
	} `json:"JDZF"`
	Jjjl struct {
		Datas []struct {
			Mgrid       string `json:"MGRID"`
			Mgrname     string `json:"MGRNAME"`
			Fcode       string `json:"FCODE"`
			Days        string `json:"DAYS"`
			Fempdate    string `json:"FEMPDATE"`
			Lempdate    string `json:"LEMPDATE"`
			Penavgrowth string `json:"PENAVGROWTH"`
			Newphotourl string `json:"NEWPHOTOURL"`
			Isinoffice  string `json:"ISINOFFICE"`
		} `json:"Datas"`
		ErrCode      int         `json:"ErrCode"`
		Success      bool        `json:"Success"`
		ErrMsg       interface{} `json:"ErrMsg"`
		Message      interface{} `json:"Message"`
		ErrorCode    string      `json:"ErrorCode"`
		ErrorMessage interface{} `json:"ErrorMessage"`
		ErrorMsgLst  interface{} `json:"ErrorMsgLst"`
		TotalCount   int         `json:"TotalCount"`
		Expansion    interface{} `json:"Expansion"`
	} `json:"JJJL"`
	Jjgm struct {
		Datas []struct {
			Fsrq   string `json:"FSRQ"`
			Netnav string `json:"NETNAV"`
			Change string `json:"CHANGE"`
			Issum  string `json:"ISSUM"`
		} `json:"Datas"`
		ErrCode      int         `json:"ErrCode"`
		Success      bool        `json:"Success"`
		ErrMsg       interface{} `json:"ErrMsg"`
		Message      interface{} `json:"Message"`
		ErrorCode    string      `json:"ErrorCode"`
		ErrorMessage interface{} `json:"ErrorMessage"`
		ErrorMsgLst  interface{} `json:"ErrorMsgLst"`
		TotalCount   int         `json:"TotalCount"`
		Expansion    string      `json:"Expansion"`
	} `json:"JJGM"`
	Hbcc struct {
		Datas struct {
			FundMMAsset struct {
				Bspctnv    string `json:"BSPCTNV"`
				Abspctnv   string `json:"ABSPCTNV"`
				Brepopctnv string `json:"BREPOPCTNV"`
				Mpctnv     string `json:"MPCTNV"`
				Oipctnv    string `json:"OIPCTNV"`
				Jzc        string `json:"JZC"`
			} `json:"FundMMAsset"`
			FundMMDistribute interface{} `json:"FundMMDistribute"`
		} `json:"Datas"`
		ErrCode      int         `json:"ErrCode"`
		Success      bool        `json:"Success"`
		ErrMsg       interface{} `json:"ErrMsg"`
		Message      interface{} `json:"Message"`
		ErrorCode    string      `json:"ErrorCode"`
		ErrorMessage interface{} `json:"ErrorMessage"`
		ErrorMsgLst  interface{} `json:"ErrorMsgLst"`
		TotalCount   int         `json:"TotalCount"`
		Expansion    string      `json:"Expansion"`
	} `json:"HBCC"`
	Fhsp struct {
		Datas struct {
			Fhinfo []struct {
				Fsrq   string `json:"FSRQ"`
				Djr    string `json:"DJR"`
				Fhfcz  string `json:"FHFCZ"`
				Cfbl   string `json:"CFBL"`
				Fhfcbz string `json:"FHFCBZ"`
				Cflx   string `json:"CFLX"`
				Ffr    string `json:"FFR"`
				Fh     string `json:"FH"`
				Dtype  string `json:"DTYPE"`
			} `json:"FHINFO"`
			Fcinfo []interface{} `json:"FCINFO"`
		} `json:"Datas"`
		ErrCode      int         `json:"ErrCode"`
		Success      bool        `json:"Success"`
		ErrMsg       interface{} `json:"ErrMsg"`
		Message      interface{} `json:"Message"`
		ErrorCode    string      `json:"ErrorCode"`
		ErrorMessage interface{} `json:"ErrorMessage"`
		ErrorMsgLst  interface{} `json:"ErrorMsgLst"`
		TotalCount   int         `json:"TotalCount"`
		Expansion    interface{} `json:"Expansion"`
	} `json:"FHSP"`
	Jjfx struct {
		Datas struct {
			Fcode          string `json:"FCODE"`
			Shortname      string `json:"SHORTNAME"`
			Fundtype       string `json:"FUNDTYPE"`
			Pltdate        string `json:"PLTDATE"`
			Bfundtype      string `json:"BFUNDTYPE"`
			Feature        string `json:"FEATURE"`
			Fsrq           string `json:"FSRQ"`
			Rzdf           string `json:"RZDF"`
			Dwjz           string `json:"DWJZ"`
			Syi            string `json:"SYI"`
			Syl            string `json:"SYL"`
			Sylname        string `json:"SYLNAME"`
			Periodname     string `json:"PERIODNAME"`
			Mcoverdetail   string `json:"MCOVERDETAIL"`
			Comethod       string `json:"COMETHOD"`
			Isbuy          string `json:"ISBUY"`
			Buy            bool   `json:"BUY"`
			Issales        string `json:"ISSALES"`
			Mindt          string `json:"MINDT"`
			Dtzt           string `json:"DTZT"`
			Appointment    string `json:"APPOINTMENT"`
			Appointmenturl string `json:"APPOINTMENTURL"`
			Shareurl       string `json:"SHAREURL"`
			Cfhid          string `json:"CFHID"`
			CFHName        string `json:"CFHName"`
			HeaderImgPath  string `json:"HeaderImgPath"`
		} `json:"Datas"`
		ErrCode    int         `json:"ErrCode"`
		ErrMsg     interface{} `json:"ErrMsg"`
		TotalCount int         `json:"TotalCount"`
		Expansion  interface{} `json:"Expansion"`
	} `json:"JJFX"`
	Jjcc struct {
		Datas struct {
			InverstPosition struct {
				FundStocks []struct {
					Gpdm         string `json:"GPDM"`
					Gpjc         string `json:"GPJC"`
					Jzbl         string `json:"JZBL"`
					Texch        string `json:"TEXCH"`
					Isinvisbl    string `json:"ISINVISBL"`
					Pctnvchgtype string `json:"PCTNVCHGTYPE"`
					Pctnvchg     string `json:"PCTNVCHG"`
					Newtexch     string `json:"NEWTEXCH"`
					Indexcode    string `json:"INDEXCODE"`
					Indexname    string `json:"INDEXNAME"`
				} `json:"fundStocks"`
				Fundboods []struct {
					Zqdm     string `json:"ZQDM"`
					Zqmc     string `json:"ZQMC"`
					Zjzbl    string `json:"ZJZBL"`
					Isbroken string `json:"ISBROKEN"`
				} `json:"fundboods"`
				Fundfofs     []interface{} `json:"fundfofs"`
				Etfcode      interface{}   `json:"ETFCODE"`
				Etfshortname interface{}   `json:"ETFSHORTNAME"`
			} `json:"InverstPosition"`
			AssetAllocation struct {
				Two0210331 []struct {
					Fsrq string `json:"FSRQ"`
					Gp   string `json:"GP"`
					Zq   string `json:"ZQ"`
					Hb   string `json:"HB"`
					Jzc  string `json:"JZC"`
					Qt   string `json:"QT"`
					Jj   string `json:"JJ"`
				} `json:"2021-03-31"`
			} `json:"AssetAllocation"`
			SectorAllocation struct {
				Two0210331 []struct {
					Hymc  string `json:"HYMC"`
					Sz    string `json:"SZ"`
					Zjzbl string `json:"ZJZBL"`
					Fsrq  string `json:"FSRQ"`
				} `json:"2021-03-31"`
			} `json:"SectorAllocation"`
		} `json:"Datas"`
		ErrCode      int         `json:"ErrCode"`
		Success      bool        `json:"Success"`
		ErrMsg       interface{} `json:"ErrMsg"`
		Message      interface{} `json:"Message"`
		ErrorCode    string      `json:"ErrorCode"`
		ErrorMessage interface{} `json:"ErrorMessage"`
		ErrorMsgLst  interface{} `json:"ErrorMsgLst"`
		TotalCount   int         `json:"TotalCount"`
		Expansion    string      `json:"Expansion"`
	} `json:"JJCC"`
	Fhtx struct {
		Datas        []interface{} `json:"Datas"`
		ErrCode      int           `json:"ErrCode"`
		Success      bool          `json:"Success"`
		ErrMsg       interface{}   `json:"ErrMsg"`
		Message      interface{}   `json:"Message"`
		ErrorCode    string        `json:"ErrorCode"`
		ErrorMessage interface{}   `json:"ErrorMessage"`
		ErrorMsgLst  interface{}   `json:"ErrorMsgLst"`
		TotalCount   int           `json:"TotalCount"`
		Expansion    interface{}   `json:"Expansion"`
	} `json:"FHTX"`
	Tssj struct {
		Datas struct {
			Sharp1         string `json:"SHARP1"`
			Sharp1Nrank    string `json:"SHARP_1NRANK"`
			Sharp1Nfsc     string `json:"SHARP_1NFSC"`
			Syl1N          string `json:"SYL_1N"`
			Maxretra1      string `json:"MAXRETRA1"`
			Maxretra1Nrank string `json:"MAXRETRA_1NRANK"`
			Maxretra1Nfsc  string `json:"MAXRETRA_1NFSC"`
			Stddev1        string `json:"STDDEV1"`
			Stddev1Nrank   string `json:"STDDEV_1NRANK"`
			Stddev1Nfsc    string `json:"STDDEV_1NFSC"`
			ProfitZ        string `json:"PROFIT_Z"`
			ProfitY        string `json:"PROFIT_Y"`
			Profit3Y       string `json:"PROFIT_3Y"`
			Profit6Y       string `json:"PROFIT_6Y"`
			Profit1N       string `json:"PROFIT_1N"`
			PvY            string `json:"PV_Y"`
			DtcountY       string `json:"DTCOUNT_Y"`
			Ffavorcount    string `json:"FFAVORCOUNT"`
			Earn1N         string `json:"EARN_1N"`
			Avghold        string `json:"AVGHOLD"`
			Brokentimes    string `json:"BROKENTIMES"`
			Isexchg        string `json:"ISEXCHG"`
			SylLn          string `json:"SYL_LN"`
			Stddev3        string `json:"STDDEV3"`
			Stddev3Nrank   string `json:"STDDEV_3NRANK"`
			Stddev3Nfsc    string `json:"STDDEV_3NFSC"`
			Stddev5        string `json:"STDDEV5"`
			Stddev5Nrank   string `json:"STDDEV_5NRANK"`
			Stddev5Nfsc    string `json:"STDDEV_5NFSC"`
			Sharp3         string `json:"SHARP3"`
			Sharp3Nrank    string `json:"SHARP_3NRANK"`
			Sharp3Nfsc     string `json:"SHARP_3NFSC"`
			Sharp5         string `json:"SHARP5"`
			Sharp5Nrank    string `json:"SHARP_5NRANK"`
			Sharp5Nfsc     string `json:"SHARP_5NFSC"`
			Maxretra3      string `json:"MAXRETRA3"`
			Maxretra3Nrank string `json:"MAXRETRA_3NRANK"`
			Maxretra3Nfsc  string `json:"MAXRETRA_3NFSC"`
			Maxretra5      string `json:"MAXRETRA5"`
			Maxretra5Nrank string `json:"MAXRETRA_5NRANK"`
			Maxretra5Nfsc  string `json:"MAXRETRA_5NFSC"`
		} `json:"Datas"`
		ErrCode      int         `json:"ErrCode"`
		Success      bool        `json:"Success"`
		ErrMsg       interface{} `json:"ErrMsg"`
		Message      interface{} `json:"Message"`
		ErrorCode    string      `json:"ErrorCode"`
		ErrorMessage interface{} `json:"ErrorMessage"`
		ErrorMsgLst  interface{} `json:"ErrorMsgLst"`
		TotalCount   int         `json:"TotalCount"`
		Expansion    interface{} `json:"Expansion"`
	} `json:"TSSJ"`
	Jjjlnew struct {
		Datas []struct {
			Days        string `json:"DAYS"`
			Fempdate    string `json:"FEMPDATE"`
			Lempdate    string `json:"LEMPDATE"`
			Penavgrowth string `json:"PENAVGROWTH"`
			Manger      []struct {
				Mgrid           string `json:"MGRID"`
				Mgrname         string `json:"MGRNAME"`
				Newphotourl     string `json:"NEWPHOTOURL"`
				Isinoffice      string `json:"ISINOFFICE"`
				Yieldse         string `json:"YIELDSE"`
				Totaldays       string `json:"TOTALDAYS"`
				Investmentidear string `json:"INVESTMENTIDEAR"`
				HjJn            int    `json:"HJ_JN"`
			} `json:"MANGER"`
		} `json:"Datas"`
		ErrCode      int         `json:"ErrCode"`
		Success      bool        `json:"Success"`
		ErrMsg       interface{} `json:"ErrMsg"`
		Message      interface{} `json:"Message"`
		ErrorCode    string      `json:"ErrorCode"`
		ErrorMessage interface{} `json:"ErrorMessage"`
		ErrorMsgLst  interface{} `json:"ErrorMsgLst"`
		TotalCount   int         `json:"TotalCount"`
		Expansion    interface{} `json:"Expansion"`
	} `json:"JJJLNEW"`
}

// QueryFundInfo 查询基金详情
func (e EastMoney) QueryFundInfo(ctx context.Context, fundCode string) (RespFundNetList, error) {
	apiurl := fmt.Sprintf("https://j5.dfcfw.com/sc/tfs/qt/v2.0.1/%v.json", fundCode)
	params := map[string]string{}
	logging.Debug(ctx, "EastMoney QueryFundInfo "+apiurl+" begin", zap.Any("params", params))
	beginTime := time.Now()
	apiurl, err := goutils.NewHTTPGetURLWithQueryString(ctx, apiurl, params)
	if err != nil {
		return RespFundNetList{}, err
	}
	resp := RespFundNetList{}
	err = goutils.HTTPGET(ctx, e.HTTPClient, apiurl, &resp)
	latency := time.Now().Sub(beginTime).Milliseconds()
	logging.Debug(
		ctx,
		"EastMoney QueryFundInfo "+apiurl+" end",
		zap.Int64("latency(ms)", latency),
		zap.Any("resp", resp),
	)
	if err != nil {
		return resp, err
	}
	if resp.ErrCode != 0 {
		return resp, fmt.Errorf("QueryFundInfo error %v", resp.ErrMsg)
	}
	return resp, nil
}
