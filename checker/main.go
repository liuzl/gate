package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/crawlerclub/x/downloader"
	"github.com/crawlerclub/x/parser"
	"github.com/crawlerclub/x/types"
	"io/ioutil"
	"log"
)

func main() {
	confFile := flag.String("conf", "./fb.json", "crawler configure file")
	proxy := flag.String("proxy", "", "proxy url")
	flag.Parse()
	conf, _ := ioutil.ReadFile(*confFile)
	var parseConf types.ParseConf
	err := json.Unmarshal(conf, &parseConf)
	if err != nil {
		log.Fatal(err)
	}

	pageUrl := "https://www.facebook.com/search/top/?q=18232032959"
	req := &types.HttpRequest{Url: pageUrl, Method: "GET", Platform: "mobile", UseProxy: false}
	if len(*proxy) > 0 {
		req = &types.HttpRequest{Url: pageUrl, Method: "GET", Platform: "mobile", UseProxy: true, Proxy: *proxy}
	}
	log.Println(*proxy)
	//fmt.Println(pageUrl)
	res := downloader.Download(req)
	if res.Error != nil {
		log.Println(res.Error)
	}
	ioutil.WriteFile("./fb.html", res.Content, 0666)
	fmt.Println(res.Text)
	_, retItems, err := parser.Parse([]byte(res.Text), pageUrl, &parseConf)
	str, _ := json.Marshal(retItems)

	fmt.Println(string(str))
}
