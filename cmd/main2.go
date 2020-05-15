package main

import (
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	sigs chan os.Signal
)

func init() {
	sigs = make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
}
func main() {

	const (
		port = 9515
	)

	//如果seleniumServer没有启动，就启动一个seleniumServer所需要的参数，可以为空，示例请参见https://github.com/tebeka/selenium/blob/master/example_test.go
	var opts []selenium.ServiceOption
	//opts := []selenium.ServiceOption{
	//    selenium.StartFrameBuffer(),           // Start an X frame buffer for the browser to run in.
	//    selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
	//}

	//selenium.SetDebug(true)
	service, err := selenium.NewChromeDriverService(`chromedriver`, port, opts...)
	if nil != err {
		fmt.Println("start a chromedriver service falid", err.Error())
		return
	}
	//注意这里，server关闭之后，chrome窗口也会关闭
	defer service.Stop()

	//链接本地的浏览器 chrome
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}

	//禁止图片加载，加快渲染速度
	imagCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}
	chromeCaps := chrome.Capabilities{
		Prefs: imagCaps,
		Path:  "",
		Args: []string{
			//"--headless", // 设置Chrome无头模式，在linux下运行，需要设置这个参数，否则会报错
			"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36", // 模拟user-agent，防反爬
		},
	}
	//以上是设置浏览器参数
	caps.AddChrome(chromeCaps)

	// 调起chrome浏览器
	driver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		fmt.Println("connect to the webDriver faild", err.Error())
		return
	}
	err = driver.Get("http://wyw.hwxnet.com/")
	if err != nil {
		panic(err)
	}
	wd, err := driver.FindElement(selenium.ByCSSSelector, "#wd")
	if err != nil {
		panic(err)
	}
	if err = wd.SendKeys("行\n"); err != nil {
		panic(err)
	}
	time.Sleep(time.Second)
	li, err := wd.FindElements(selenium.ByCSSSelector, "#content > div.sub_con.f18.clearfix > ul")
	if err == nil {
		for i := range li {
			pinyin, _ := li[i].FindElement(selenium.ByCSSSelector, ".pinyin")
			fmt.Println(pinyin)
			a, _ := li[i].FindElement(selenium.ByCSSSelector, "a")
			href, _ := a.GetAttribute("href")
			fmt.Println(a, href)
		}
	}
	<-make(chan struct{})
}
