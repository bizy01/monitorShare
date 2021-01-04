// his_flow_url = 'http://push2his.eastmoney.com/api/qt/kamt.kline/get?fields1=f1,f3,f5&fields2=f51,f52&klt=101&lmt=300'
// realtime_flow_url = 'http://push2.eastmoney.com/api/qt/kamt.rtmin/get?fields1=f1,f3&fields2=f51,f52,f54,f56'

package spider

import (
	"net/http"
	"regexp"
	"errors"
	"github.com/tidwall/gjson"
	"github.com/bizy01/monitorShare/cliutils"
)


const (
	northshUrl ="http://dcfm.eastmoney.com/EM_MutiSvcExpandInterface/api/js/get?type=HSGTHIS&token=70f12f2f4f091e459a279469fe49eca5&filter=(MarketType=1)&js=var%20FvPygxqC={%22data%22:(x),%22pages%22:(tp)}&ps=1&p=1&sr=-1&st=DetailDate&rt=53658974"
    northszUrl ="http://dcfm.eastmoney.com/EM_MutiSvcExpandInterface/api/js/get?type=HSGTHIS&token=70f12f2f4f091e459a279469fe49eca5&filter=(MarketType=3)&js=var%20FvPygxqC={%22data%22:(x),%22pages%22:(tp)}&ps=1&p=1&sr=-1&st=DetailDate&rt=53658974"
    northRegMatchPattern = "=(.*?)$"
)
type NorthFund struct {
    httpCli *http.Client
    MarkType map[string]string
    FundData *FundData
}

type FundData struct {
	MarketType string  // 市场类型
	Datetime       string  // 时间
	TotalNetIn   float64  // 当日成交净买额
	TotalIn float64  // 买入额
	TotalOut   float64  // 卖出额
	TodayIn  float64   // 当日资金流入
	GrandTotalIn  float64   // 历史累计流入
	TodayBalance    float64   // 当日余额
	StockCode     string   // 领涨股代码
	Stock     string   // 领涨股名称
	StockUp     float64   // 领涨股涨跌幅
	Index     float64   // 指数
	IndexPercent     float64   // 指数涨跌幅
}

func NewNorthFund() *NorthFund {
	httpClient := &http.Client{}

	return &NorthFund{
		httpCli: httpClient,
		MarkType: map[string]string{
			"sh": "1",
			"sz": "3",
		},
		FundData: new(FundData),
	}
}

// 获取指数数据
func (s *NorthFund) GetData() {
	var statuCode int
	var resp interface{}

	for name, code := range s.MarkType {
		if code == "1" {
    		statuCode, resp = cliutils.HttpCli(s.httpCli, "get", northshUrl, "", nil)
		} else {
    		statuCode, resp = cliutils.HttpCli(s.httpCli, "get", northszUrl, "", nil)
		}

    	if statuCode / 100 == 2 {
    		err := s.parse(name, resp.(string))
    		if err != nil {

    		}
    	}
	}
}

func (s *NorthFund) parse(code, buf string) error {
	reg := regexp.MustCompile(northRegMatchPattern)
	if reg == nil {
		return errors.New("MustCompile err")
	}

	//提取关键信息
	result := reg.FindStringSubmatch(buf)[1]

	for _, item := range gjson.Get(result, "data").Array() {
		s.FundData.MarketType = item.Get("MarketType").String()
		s.FundData.Datetime = item.Get("DetailDate").String()
		s.FundData.TotalNetIn = item.Get("DRCJJME").Float()
		s.FundData.TotalIn = item.Get("MRCJE").Float()
		s.FundData.TotalOut = item.Get("MCCJE").Float()
		s.FundData.TodayIn = item.Get("DRZJLR").Float()
		s.FundData.GrandTotalIn = item.Get("LSZJLR").Float()
		s.FundData.TodayBalance  = item.Get("DRYE").Float()
		s.FundData.StockCode = item.Get("LCGCode").String()
		s.FundData.Stock = item.Get("LCG").String()
		s.FundData.StockUp = item.Get("LCGZDF").Float()
		s.FundData.Index = item.Get("SSEChange").Float()
		s.FundData.IndexPercent = item.Get("SSEChangePrecent").Float()
	}

	return nil
}

// 渲染报告
func (s *NorthFund) renderReport() {

}

// 发送报告
func (s *NorthFund) SendReport() {

}