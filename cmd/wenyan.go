package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/morikuni/aec"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
)

var key = "行"
var pinyin = ""
var url = false

func main() {
	flag.StringVar(&key, "k", "", "-k key:输入需要翻译文言文词语")
	flag.StringVar(&pinyin, "p", "", "-p pinyin:输入需要翻译的文字拼音")
	flag.BoolVar(&url, "u", false, "-u url:打开页面")
	flag.Parse()
	if key == "" && pinyin == "" {
		notfound(key)
	}
	if key != "" {
		keyPrint(key)
	} else if pinyin != "" {
		pinyinPrint(pinyin)
	}
}

func pinyinPrint(pinyin string) {
	dir, _ := os.UserCacheDir()
	file, err := ioutil.ReadFile(path.Join(dir, ".wenyan", "pcache", pinyin+".json"))
	if err != nil {
		notfound(pinyin)
	}
	pinyins := make([][2]string, 0)
	if err = json.Unmarshal(file, &pinyins); err != nil {
		notfound(pinyin)
	}

	for i, strings := range pinyins {
		fmt.Printf(aec.RedF.Apply("%d")+":%s  ", i, strings)
		if i%3 == 0 && i != 0 {
			fmt.Println()
		}
	}
in:
	fmt.Println("请输入编号:")
	in := bufio.NewReader(os.Stdin)
	idx, _, _ := in.ReadLine()

	fmt.Println(aec.BlackF.Apply("-----------"))
	i, err := strconv.Atoi(string(idx))
	if err != nil || i >= len(pinyins) {
		goto in
	}
	keyPrint(pinyins[i][1])
}
func keyPrint(key string) {
	dir, _ := os.UserCacheDir()
	key = string([]rune(key)[0])
	file, err := ioutil.ReadFile(path.Join(dir, ".wenyan", "cache", key+".json"))
	if err != nil {
		notfound(key)
	}
	m := WenYan{}
	err = json.Unmarshal(file, &m)
	if err != nil {
		panic(err)
	}
	fmt.Println(m)
}
func notfound(key string) {
	fmt.Printf("没有找到与您查询的“%s”相关的结果。\n", key)
	os.Exit(1)
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
	s := ""
	if url {
		s = fmt.Sprintf(`%s 拼音:%v 部首:%s 部首笔画:%s 总笔画:%s 笔顺:%s
详细释义:
%s
与“%s”相关的词语:
%v
与“%s”相关的成语:
%v
链接:
%s`, aec.Bold.Apply(w.Key), w.Pinyin, w.BuShou, w.BuShouBiHua, w.TotalBiHua, w.BiShun, w.Desc, w.Key, w.word, w.Key, w.CY, w.URL)
	} else {
		s = fmt.Sprintf(`%s 拼音:%v 部首:%s 部首笔画:%s 总笔画:%s 笔顺:%s
详细释义:
%s
与“%s”相关的词语:
%v
与“%s”相关的成语:
%v`, aec.Bold.Apply(w.Key), w.Pinyin, w.BuShou, w.BuShouBiHua, w.TotalBiHua, w.BiShun, w.Desc, w.Key, w.word, w.Key, w.CY)
	}

	reg, err := regexp.Compile(w.Key)
	if err != nil {
		panic(err)
	}
	reg2, err := regexp.Compile("<")
	if err != nil {
		panic(err)
	}
	return reg2.ReplaceAllString(reg.ReplaceAllString(s, aec.RedF.Apply(w.Key)), " <")
}
