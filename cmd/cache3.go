package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	get, err := http.Get("http://wyw.hwxnet.com/pinyin.html")
	if err != nil {
		panic(err)
	}
	doc, err := goquery.NewDocumentFromReader(get.Body)
	if err != nil {
		panic(err)
	}
	var urlss = map[string]string{}
	doc.Find("#content > div:nth-child(4)").Find("a").Each(func(i int, s *goquery.Selection) {
		attr, ok := s.Attr("href")
		if ok {
			space := strings.TrimSpace(s.Text())
			url := "http://wyw.hwxnet.com" + attr
			fmt.Println(space, url)
			resp, err := http.Get(url)
			if err != nil {
				panic(err)
			}
			doc2, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				panic(err)
			}
			doc2.Find("#content > div:nth-child(6) > dl > dd > a").Each(func(i int, s *goquery.Selection) {
				val, ok := s.Attr("href")
				if ok {
					urlss[strings.TrimSpace(s.Text())] = "http://wyw.hwxnet.com" + val
				}
			})
		}
	})
	marshal, err := json.Marshal(urlss)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile("cgo/pinyin.json", marshal, 0777)
}
