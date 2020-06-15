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
	file, err := ioutil.ReadFile("cgo/pinyin.json")
	if err != nil {
		panic(err)
	}
	var s = map[string]string{}
	if err = json.Unmarshal(file, &s); err != nil {
		panic(err)
	}
	count := 1
	lenth := len(s)
	for k, v := range s {
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
