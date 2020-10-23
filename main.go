package main

import (
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/parnurzeal/gorequest"
	"golang.org/x/net/html"
	"log"
	"regexp"
	"strconv"
)
//建立行政区域数据结构体
type adminRegion struct{
	RegionCode int
	Name string
	Provincecode int
	Citycode int
	Countycode int
}
//利用gorequest抓取网页，htmlquery解析网页（传入网址，返回html节点）
func fetch(url string) *html.Node {
	request := gorequest.New()
	resp, _, _ := request.Get(url).
		Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36").
		End()

	doc, err := htmlquery.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

//分析提取数据（传入需解析网址，无返回值）
func parseUrls(url string) {
	doc := fetch(url)
	nodes := htmlquery.Find(doc, `//*[@id="2020年7月份县以上行政区划代码_24734"]/table/tbody//tr`)
	var province [83] string
	var city [66][100] string
	for _, node := range nodes {
		//以下部分代码分析获取省市的代码和名字
		provinceCityInfo := htmlquery.Find(node, `//td[@class="xl7024734"]`)

		if provinceCityInfo != nil{
			regionCodePc := htmlquery.InnerText(htmlquery.FindOne(provinceCityInfo[0],`./text()`))
			match,_ := regexp.MatchString("[1-8]{2}0000", regionCodePc)//判断代码表示的省还是市

			if match {
				//获取省的代码和名字
				provinceNum,_ := strconv.Atoi(regionCodePc[0:2])
				provinceName := htmlquery.InnerText(htmlquery.FindOne(provinceCityInfo[1],`./text()`))
				province[provinceNum] = provinceName
			}else{
				//获取市的代码和名字
				provinceNum,_ := strconv.Atoi(regionCodePc[0:2])
				cityNum,_ := strconv.Atoi(regionCodePc[2:4])
				cityName := htmlquery.InnerText(htmlquery.FindOne(provinceCityInfo[1],`./text()`))
				city[provinceNum][cityNum] = cityName
			}
		}

		//以下部分代码分析获取县的代码和名字
		placeInfos := htmlquery.Find(node, `//td[@class="xl7124734"]`)

		if placeInfos != nil {
			//获取县的代码和名字
			regionCode, _ := strconv.Atoi(htmlquery.InnerText(htmlquery.FindOne(placeInfos[0],`.//text()`)))
			locationName := htmlquery.InnerText(htmlquery.FindOne(placeInfos[1],`./text()`))

			//获取省和市名字数组索引号
			provinceNum := regionCode / 10000
			cityNum := regionCode / 100 % 100

			//合成总地区名，生成JSON数据
			name := province[provinceNum] + city[provinceNum][cityNum] + locationName
			toJson(regionCode,name,provinceNum,cityNum)
		}
	}
}

//合成json并输出（传入县级行政代码，县级名，省代码，市代码,无返回值）
func toJson(regionCode int,name string,provinceNum int,cityNum int){

	aR := &adminRegion{
		regionCode,
		name,
		provinceNum * 10000,
		(provinceNum * 100 + cityNum) * 100 ,
		regionCode,
	}

	aRJ, err := json.Marshal(aR)
	if err != nil {
		fmt.Println("encoding failed")
	} else {
		fmt.Println("  ")
		//fmt.Println(aRJ)
		fmt.Println(string(aRJ))
	}
}

func main(){
	parseUrls("http://www.mca.gov.cn//article/sj/xzqh/2020/2020/2020092500801.html")
}
