package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/morikuni/aec"
	"golang.org/x/net/html"
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
		cache1(urls[i])
	}
	cache2()
}
func cache1(url string) {
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
func cache2() {
	count := 1
	for s := range index {
		get, err := http.Get(s)
		if err != nil {
			panic(err)
		}
		doc, err := goquery.NewDocumentFromReader(get.Body)
		if err != nil {
			panic(err)
		}
		yan := WenYan{URL: s, Key: index[s]}
		doc.Find("#content > div.word_con.clearfix > div.introduce > div:nth-child(1) .pinyin").Each(func(i int, s *goquery.Selection) {
			yan.Pinyin = append(yan.Pinyin, strings.TrimSpace(s.Text()))
		})
		bushou := doc.Find("#content > div.word_con.clearfix > div.introduce > div:nth-child(2) > span:nth-child(2)")
		yan.BuShou = strings.TrimSpace(bushou.Text())
		bushouSize := doc.Find("#content > div.word_con.clearfix > div.introduce > div:nth-child(2) > span:nth-child(4)")
		yan.BuShouBiHua = strings.TrimSpace(bushouSize.Text())
		totalSize := doc.Find("#content > div.word_con.clearfix > div.introduce > div:nth-child(2) > span:nth-child(6)")
		yan.TotalBiHua = strings.TrimSpace(totalSize.Text())
		bishun := doc.Find("#content > div.word_con.clearfix > div.introduce > div:nth-child(3) > span")
		yan.BiShun = strings.TrimSpace(bishun.Text())
		yan.Desc = strings.TrimSpace(Text(doc.Find("#content > div.view_con.clearfix")))
		yan.word = Str2Arr(doc.Find("#content > div:nth-child(10) > ul").Text())
		yan.CY = Str2Arr(doc.Find("#content > div:nth-child(12) > ul").Text())
		marshal, err := json.Marshal(&yan)
		if err != nil {
			panic(err)
		}
		filename := "cgo/cache/" + yan.Key + ".json"
		_, err = os.Stat(filename)
		if err == nil {
			fmt.Println(count, "/", len(index))
			count++
			continue
		} else {
			err = ioutil.WriteFile(filename, marshal, 0777)
			if err != nil {
				panic(err)
			}
			fmt.Println(count, "/", len(index))
			count++
		}
	}
}
func Str2Arr(str string) (arr []string) {
	split := strings.Split(str, "\n")
	for i := range split {
		space := strings.TrimSpace(split[i])
		if space != "" {
			arr = append(arr, space)
		}
	}
	return arr
}
func Text(s *goquery.Selection) string {
	var buf bytes.Buffer
	// Slightly optimized vs calling Each: no single selection object created
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Data == "br" {
			buf.WriteString("\n")
		}
		if n.Type == html.TextNode {
			// Keep newlines and spaces, like jQuery
			buf.WriteString(n.Data)
		}
		if n.FirstChild != nil {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}
	for _, n := range s.Nodes {
		f(n)
	}

	return buf.String()
}

type WenYan struct {
	Key         string
	URL         string
	Pinyin      []string
	BuShou      string
	BuShouBiHua string
	TotalBiHua  string
	BiShun      string
	Desc        string
	word        []string
	CY          []string
}

func (w WenYan) String() string {
	s := fmt.Sprintf(`%s 拼音:%v 部首:%s 部首笔画:%s 总笔画:%s 笔顺:%s
详细释义:
%s
与“%s”相关的词语:
%v
与“%s”相关的成语:
%v
链接:
%s`, w.Key, w.Pinyin, w.BuShou, w.BuShouBiHua, w.TotalBiHua, w.BiShun, w.Desc, w.Key, w.word, w.Key, w.CY, w.URL)
	reg, err := regexp.Compile(w.Key)
	if err != nil {
		panic(err)
	}
	return reg.ReplaceAllString(s, aec.RedF.Apply(w.Key))
}
