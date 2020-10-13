package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var urlss = map[string]string{}

func main() {
	load()
	cache4()
}

func cache4() {
	count := 1
	lenth := len(urlss)
	for k, v := range urlss {
		save := make([][2]string, 0)
		get, err := http.Get(v)
		if err != nil {
			panic(err)
		}
		doc, err := goquery.NewDocumentFromReader(get.Body)
		if err != nil {
			panic(err)
		}
		doc.Find("#content > div.sub_con.f14.clearfix > ul a").Each(func(i int, se *goquery.Selection) {
			py := strings.TrimSpace(se.Find(".py").Text())
			zi := strings.TrimSpace(strings.TrimLeft(strings.TrimSpace(se.Text()), py))
			save = append(save, [2]string{py, zi})
			fmt.Println(py, zi)
		})

		marshal, err := json.Marshal(save)
		if err != nil {
			panic(err)
		}
		ioutil.WriteFile("cgo/pcache/"+k+".json", marshal, 0777)
		fmt.Println(count, "/", lenth)
		count++
	}
}

func load() {
	get, err := http.Get("http://wyw.hwxnet.com/pinyin.html")
	if err != nil {
		panic(err)
	}
	doc, err := goquery.NewDocumentFromReader(get.Body)
	if err != nil {
		panic(err)
	}
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
}
