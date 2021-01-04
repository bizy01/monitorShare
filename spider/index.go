package spider

import (
	"net/http"
	"fmt"
	"regexp"
	"errors"
	"strings"
	"strconv"
	"github.com/tidwall/gjson"
	"github.com/bizy01/monitorShare/cliutils"
)


const (
	indexPattern ="http://push2his.eastmoney.com/api/qt/stock/kline/get?cb=jQuery1124016963078166077672_1576074523699&secid=%s&fields1=f1,f2,f3,f4,f5&fields2=f51,f52,f53,f54,f55,f56,f57,f58&klt=101&fqt=0&beg=%s&end=%s"
    indexRegMatchPattern = `\((.*?)\)`
)
// http://push2his.eastmoney.com/api/qt/stock/kline/get?cb=jQuery1124016963078166077672_1576074523699&secid=1.000001&fields1=f1%2Cf2%2Cf3%2Cf4%2Cf5&fields2=f51%2Cf52%2Cf53%2Cf54%2Cf55%2Cf56%2Cf57%2Cf58&klt=101&fqt=0&beg=19900101&end=20300101
type StockIndex struct {
    httpCli *http.Client
    indexCode map[string]string
    IndexData map[string]*IndexData
}

type IndexData struct {
	Code       string  // code
	Name       string  // name
	Datetime   string  // 时间
	PriceStart float64  // 开盘价
	PriceEnd   float64  // 收盘价
	PriceMax  float64   // 最高价
	PriceMin  float64   // 最低价
	Amount    float64   // 成交量
	Value     float64   // 成交额
	Swing     float64   // 振幅
}

func NewStockIndex() *StockIndex {
	httpClient := &http.Client{}

	return &StockIndex{
		httpCli: httpClient,
		indexCode: map[string]string{
			"sh": "1.000001",
			"sz": "0.399001",
			"cj": "0.399006",
		},
		IndexData: make(map[string]*IndexData),
	}
}

// 获取指数数据
func (s *StockIndex) GetIndex() {
	// 上证指数
	for name, code := range s.indexCode {
		path := fmt.Sprintf(indexPattern, code, "20201230", "20210103")
    	statuCode, resp := cliutils.HttpCli(s.httpCli, "get", path, "", nil)
    	if statuCode / 100 == 2 {
    		err := s.parse(name, resp.(string))
    		if err != nil {

    		}
    	}
	}
}

func (s *StockIndex) parse(code, buf string) error {
	reg := regexp.MustCompile(indexRegMatchPattern)
	if reg == nil {
		return errors.New("MustCompile err")
	}

	//提取关键信息
	result := reg.FindStringSubmatch(buf)[1]

	s.IndexData = map[string]*IndexData{
		code: new(IndexData),
	}

	s.IndexData[code].Code = gjson.Get(result, "data.code").String()
	s.IndexData[code].Name = gjson.Get(result, "data.name").String()

	for _, item := range gjson.Get(result, "data.klines").Array() {
		itemArr := strings.Split(item.String(), ",")

		s.IndexData[code].Datetime = itemArr[0]
		s.IndexData[code].PriceStart, _ = strconv.ParseFloat(itemArr[1], 64)
		s.IndexData[code].PriceEnd, _ = strconv.ParseFloat(itemArr[2], 64)
		s.IndexData[code].PriceMax, _ = strconv.ParseFloat(itemArr[3], 64)
		s.IndexData[code].PriceMin, _ = strconv.ParseFloat(itemArr[4], 64)
		s.IndexData[code].Amount, _ = strconv.ParseFloat(itemArr[5], 64)
		s.IndexData[code].Value, _  = strconv.ParseFloat(itemArr[6], 64)
		s.IndexData[code].Swing, _ = strconv.ParseFloat(itemArr[7], 64)
	}

	fmt.Printf("code = %s \n", code)
	fmt.Printf("value = %v \n", s.IndexData[code])

	return nil
}

// 渲染报告
func (s *StockIndex) renderReport() {

}

// 发送报告
func (s *StockIndex) SendReport() {

}