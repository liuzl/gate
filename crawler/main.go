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
	confFile := flag.String("conf", "./free-proxy-list.json", "crawler configure file")
	flag.Parse()
	conf, _ := ioutil.ReadFile(*confFile)
	var parseConf types.ParseConf
	err := json.Unmarshal(conf, &parseConf)
	if err != nil {
		log.Fatal(err)
	}

	pageUrl := "https://free-proxy-list.net/"
	req := &types.HttpRequest{Url: pageUrl, Method: "GET", Platform: "pc", UseProxy: false}
	//fmt.Println(pageUrl)
	res := downloader.Download(req)
	if res.Error != nil {
		log.Println(res.Error)
	}
	_, retItems, err := parser.Parse([]byte(res.Text), pageUrl, &parseConf)
	str, _ := json.Marshal(retItems)

	fmt.Println(string(str))
}
