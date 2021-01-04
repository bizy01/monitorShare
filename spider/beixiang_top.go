// his_flow_url = 'http://push2his.eastmoney.com/api/qt/kamt.kline/get?fields1=f1,f3,f5&fields2=f51,f52&klt=101&lmt=300'
// realtime_flow_url = 'http://push2.eastmoney.com/api/qt/kamt.rtmin/get?fields1=f1,f3&fields2=f51,f52,f54,f56'

package spider

import (
	"net/http"
	"fmt"
	"regexp"
	"errors"
	"github.com/tidwall/gjson"
	"github.com/bizy01/monitorShare/cliutils"
)


const (
	northTopUrl ="http://data.eastmoney.com/hsgt/top10.html"
    northSHPattern = "var DATA1 = (.*?);"
    northZHPattern = "var DATA2 = (.*?);"
)
type NorthTop struct {
    httpCli *http.Client
    SHTopData []*TopData
    SZTopData []*TopData
}

type TopData struct {
	MarketType string  // 市场类型
	Datetime   string  // 时间
	Rank       string  // rank
	Code       string  // 股码
	Name       string  // 名称
	Close      float64   // 当日收盘价
	ChangePercent  float64   // 当日涨跌幅
	NetIn     float64   // 当日净流入
	In        float64   // 当日流入
	Out       float64   // 当日流出
	Total     float64   // 当日成交额
}

func NewNorthTop() *NorthTop {
	httpClient := &http.Client{}

	return &NorthTop{
		httpCli: httpClient,
		SHTopData: make([]*TopData,0),
		SZTopData: make([]*TopData,0),
	}
}

// 获取指数数据
func (s *NorthTop) GetData() {
    statuCode, resp := cliutils.HttpCli(s.httpCli, "get", northTopUrl, "", nil)
	if statuCode / 100 == 2 {
		err := s.parse("", resp.(string))
		if err != nil {

		}
	}
}

func (s *NorthTop) parse(code, buf string) error {
	reg := regexp.MustCompile(northSHPattern)
	if reg == nil {
		return errors.New("MustCompile err")
	}

	//提取关键信息
	result := reg.FindStringSubmatch(buf)[1]

	for _, item := range gjson.Get(result, "data").Array() {
		topData := new(TopData)

		topData.MarketType = item.Get("MarketType").String()
		topData.Datetime = item.Get("DetailDate").String()
		topData.Rank = item.Get("Rank").String()
		topData.Code = item.Get("Code").String()
		topData.Name = item.Get("Name").String()
		topData.Close = item.Get("Close").Float()
		topData.ChangePercent = item.Get("ChangePercent").Float()
		topData.NetIn  = item.Get("HGTJME").Float()
		topData.In = item.Get("HGTMRJE").Float()
		topData.Out = item.Get("HGTMCJE").Float()
		topData.Total = item.Get("HGTCJJE").Float()

		s.SHTopData = append(s.SHTopData, topData)
	}

	reg = regexp.MustCompile(northZHPattern)
	if reg == nil {
		return errors.New("MustCompile err")
	}

	//提取关键信息
	result = reg.FindStringSubmatch(buf)[1]

	for _, item := range gjson.Get(result, "data").Array() {
		topData := new(TopData)

		topData.MarketType = item.Get("MarketType").String()
		topData.Datetime = item.Get("DetailDate").String()
		topData.Rank = item.Get("Rank").String()
		topData.Code = item.Get("Code").String()
		topData.Name = item.Get("Name").String()

		fmt.Println("+++++++", topData.Name)
		topData.Close = item.Get("Close").Float()
		topData.ChangePercent = item.Get("ChangePercent").Float()
		topData.NetIn  = item.Get("HGTJME").Float()
		topData.In = item.Get("HGTMRJE").Float()
		topData.Out = item.Get("HGTMCJE").Float()
		topData.Total = item.Get("HGTCJJE").Float()

		s.SZTopData = append(s.SZTopData, topData)
	}

	fmt.Println("sh =====>", s.SHTopData[0])
	fmt.Println("sz =====>", s.SZTopData[0])

	return nil
}

// 渲染报告
func (s *NorthTop) renderReport() {

}

// 发送报告
func (s *NorthTop) SendReport() {

}