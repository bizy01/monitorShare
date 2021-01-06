// his_flow_url = 'http://push2his.eastmoney.com/api/qt/kamt.kline/get?fields1=f1,f3,f5&fields2=f51,f52&klt=101&lmt=300'
// realtime_flow_url = 'http://push2.eastmoney.com/api/qt/kamt.rtmin/get?fields1=f1,f3&fields2=f51,f52,f54,f56'

package spider

import (
	"net/http"
	"regexp"
	"errors"
	"fmt"
	"encoding/json"
	"strings"
	"strconv"
	// "github.com/tidwall/gjson"
	"github.com/bizy01/monitorShare/cliutils"
)
// http://fund.eastmoney.com/data/rankhandler.aspx?op=ph&dt=kf&ft=all&rs=&gs=0&sc=zzf&st=desc&sd=2020-01-06&ed=2021-01-06&qdii=&tabSubtype=,,,,,&pi=1&pn=50&dx=1&v=0.4740296280780578

const (
	foundUrl ="http://fund.eastmoney.com/data/rankhandler.aspx?op=ph&dt=kf&ft=all&rs=&gs=0&sc=zzf&st=desc&sd=2020-01-06&ed=2021-01-06&qdii=&tabSubtype=,,,,,&pi=1&pn=50&dx=1&v=0.4740296280780578"
    foundRegMatchPattern = "datas:[(.*?)],"
    cookie = `xsb_history=833153%7C%u5267%u661F%u4F20%u5A92; st_si=05268533674488; qgqp_b_id=c7c6fe3868158138077d3e06f54fa5e5; intellpositionL=1140px; ASP.NET_SessionId=knoabg2pvjxk2frell40d42i; cowCookie=true; EMFUND0=01-06%2021%3A44%3A04@%23%24%u62DB%u5546%u884C%u4E1A%u7CBE%u9009%u80A1%u7968%u57FA%u91D1@%23%24000746; EMFUND1=01-06%2021%3A44%3A19@%23%24%u5609%u5B9E%u6D88%u8D39%u7CBE%u9009%u80A1%u7968A@%23%24006604; EMFUND2=01-06%2021%3A44%3A26@%23%24%u5609%u5B9E%u73AF%u4FDD%u4F4E%u78B3%u80A1%u7968@%23%24001616; EMFUND3=01-06%2021%3A44%3A36@%23%24%u56FD%u6CF0%u667A%u80FD%u6C7D%u8F66%u80A1%u7968@%23%24001790; EMFUND4=01-06%2021%3A44%3A47@%23%24%u6C47%u6DFB%u5BCC%u4E2D%u8BC1%u65B0%u80FD%u6E90%u6C7D%u8F66A@%23%24501057; EMFUND5=01-06%2021%3A44%3A55@%23%24%u5609%u5B9E%u6D88%u8D39%u7CBE%u9009%u80A1%u7968C@%23%24006605; EMFUND6=01-06%2021%3A51%3A08@%23%24%u519C%u94F6%u5DE5%u4E1A4.0%u6DF7%u5408@%23%24001606; EMFUND7=01-06%2021%3A50%3A34@%23%24%u519C%u94F6%u65B0%u80FD%u6E90%u4E3B%u9898@%23%24002190; EMFUND8=01-06%2021%3A50%3A51@%23%24%u519C%u94F6%u7814%u7A76%u7CBE%u9009%u6DF7%u5408@%23%24000336; EMFUND9=01-06 21:52:20@#$%u534E%u590F%u6D88%u8D39ETF@%23%24510630; waptgshowtime=202116; st_asi=delete; intellpositionT=1286px; st_pvi=71370113789103; st_sp=2020-11-03%2018%3A33%3A58; st_inirUrl=https%3A%2F%2Fwww.google.com%2F; st_sn=106; st_psi=20210106221312545-0-2625925376`
)
type FundRank struct {
    httpCli *http.Client
    FundData []*Data
}

type Data struct {
	Code string  // 市场类型
	Name string  // 时间
	Datetime string
	DayGrow   float64  // 当日成交净买额
	WeekGrow  float64  // 买入额
	MonthGrow float64  // 卖出额
	ThreeMonthGrow   float64   // 当日资金流入
	SixMonthGrow  float64   // 历史累计流入
	YearGrow    float64   // 当日余额
	TwoYearGrow    float64   // 当日余额
	ThreeYearGrow    float64   // 当日余额
	CurrentGrow float64
	BuildGrow float64
}

func NewFundRank() *FundRank {
	httpClient := &http.Client{}

	return &FundRank{
		httpCli: httpClient,
		FundData: make([]*Data, 0),
	}
}

// 获取指数数据
func (s *FundRank) GetData() {
	var statuCode int
	var resp interface{}

	headers := map[string]string{
		"Cookie":           cookie,
		"Host":             "fund.eastmoney.com",
		"Connection":       "keep-alive",
		"Referer": "http://fund.eastmoney.com/data/fundranking.html",
		"User-Agent":       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36",
	}

    statuCode, resp = cliutils.HttpCli(s.httpCli, "get", foundUrl, "", headers)

	if statuCode / 100 == 2 {
		err := s.parse("", resp.(string))
		if err != nil {

		}
	}
}

func (s *FundRank) parse(code, buf string) error {
	reg := regexp.MustCompile(foundRegMatchPattern)
	if reg == nil {
		return errors.New("MustCompile err")
	}

	//提取关键信息
	result := reg.FindStringSubmatch(buf)[1]

	for _, item := range data["datas"] {
		line := item

		rankData := new(Data)
		// 15.12,28.45,35.26,79.95,99.45,140.33,86.41,9.65,45.40,2015-12-02,1,99.4513,1.50%,0.15%,1,0.15%,1,46.87
		arr := strings.Split(line, ",")

		rankData.Code = arr[0]
		rankData.Name = arr[1]
		rankData.Datetime = arr[3]
		rankData.DayGrow, _ = strconv.ParseFloat(arr[6], 64)
		rankData.WeekGrow, _ = strconv.ParseFloat(arr[7], 64)
		rankData.MonthGrow, _ = strconv.ParseFloat(arr[8], 64)
		rankData.ThreeMonthGrow, _ = strconv.ParseFloat(arr[9], 64)
		rankData.SixMonthGrow, _ = strconv.ParseFloat(arr[10], 64)
		rankData.YearGrow, _ = strconv.ParseFloat(arr[11], 64)
		rankData.TwoYearGrow, _ = strconv.ParseFloat(arr[12], 64)
		rankData.ThreeYearGrow, _ = strconv.ParseFloat(arr[13], 64)
		rankData.CurrentGrow, _ = strconv.ParseFloat(arr[14], 64)
		rankData.BuildGrow, _ = strconv.ParseFloat(arr[15], 64)

		fmt.Println("result ====>", rankData)

		s.FundData = append(s.FundData, rankData)
	}

	return nil
}

// 渲染报告
func (s *FundRank) renderReport() {

}

// 发送报告
func (s *FundRank) SendReport() {

}