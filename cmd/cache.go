package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var urls = []string{
	"http://wyw.hwxnet.com/bushou/hwxE4hwxB8hwx80_1.html",
	"http://wyw.hwxnet.com/bushou/hwxE4hwxB8hwxB7_2.html",
	"http://wyw.hwxnet.com/bushou/hwxE5hwx8FhwxA3_3.html",
	"http://wyw.hwxnet.com/bushou/hwxE5hwx8Ehwx84_4.html",
	"http://wyw.hwxnet.com/bushou/hwxE6hwxAFhwx8D_5.html",
	"http://wyw.hwxnet.com/bushou/hwxE7hwxABhwxB9_6.html",
	"http://wyw.hwxnet.com/bushou/hwxE8hwxA7hwx92_7.html",
	"http://wyw.hwxnet.com/bushou/hwxE9hwx87hwx91_8.html",
	"http://wyw.hwxnet.com/bushou/hwxE9hwx9DhwxA2_9.html",
	"http://wyw.hwxnet.com/bushou/hwxE9hwxA6hwxAC_10.html",
	"http://wyw.hwxnet.com/bushou/hwxE9hwxADhwx9A_11.html",
	"http://wyw.hwxnet.com/bushou/hwxE9hwxBBhwx8D_12.html",
	"http://wyw.hwxnet.com/bushou/hwxE9hwxBChwx93_13.html",
	"http://wyw.hwxnet.com/bushou/hwxE9hwxBDhwx92_15.html",
	"http://wyw.hwxnet.com/bushou/hwxE9hwxBEhwxA0_17.html",
}

var index = map[string]string{}

func main() {
	for i := range urls {
		ss(urls[i])
		fmt.Printf("%d/%d\n", i+1, len(urls))
	}
	marshal, err := json.Marshal(index)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile("index.json", marshal, 0777)
}
func ss(url string) {
	get, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	doc, err := goquery.NewDocumentFromReader(get.Body)
	if err != nil {
		panic(err)
	}
	doc.Find("#content > div:nth-child(6) a ").Each(func(i int, s *goquery.Selection) {
		attr, ok := s.Attr("href")
		if ok {
			//fmt.Println(strings.TrimSpace(s.Text()), attr)
			resp, err := http.Get("http://wyw.hwxnet.com/" + attr)
			if err != nil {
				panic(err)
			}
			doc2, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				panic(err)
			}
			doc2.Find("#content > div.sub_con.f14.clearfix a ").Each(func(i int, s2 *goquery.Selection) {
				attr, ok := s2.Attr("href")
				if ok {
					index[attr] = strings.TrimSpace(s2.Text())
					fmt.Println(s2.Text(), attr)
				}
			})
		}
	})
}
