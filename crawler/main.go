package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/crawlerclub/x/downloader"
	"github.com/crawlerclub/x/parser"
	"github.com/crawlerclub/x/types"
	"io/ioutil"
	"log"
	"os"
)

func check(proxy string) string {
	pageUrl := "https://www.google.com"
	req := &types.HttpRequest{Url: pageUrl, Method: "GET", Platform: "mobile", UseProxy: true, Proxy: proxy, Timeout: 5}
	res := downloader.Download(req)
	if res.Error != nil {
		log.Println(res.Error)
		return "FAIL"
	} else {
		return "SUCC"
	}
}

func main() {
	confFile := flag.String("conf", "./free-proxy-list.json", "crawler configure file")
	outFile := flag.String("out", "out.txt", "the result file")
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
	//str, _ := json.Marshal(retItems)
	//fmt.Println(string(str))
	f, err := os.Create(*outFile)
	if err != nil {
		log.Println("create file error:", *outFile)
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	if len(retItems) > 0 {
		proxies := retItems[0]["proxies"]
		for _, p := range proxies.([]interface{}) {
			proxy := p.(map[string]interface{})
			if ip, ok := proxy["Ip"]; ok {
				schema := "http://"
				if proxy["Https"].(string) == "yes" {
					schema = "https://"
				}
				one := fmt.Sprintf("%s%s:%s", schema, ip, proxy["Port"])
				line := fmt.Sprintf("%s\t%s", one, check(one))
				fmt.Println(line)
				fmt.Fprintln(w, line)
				w.Flush()
			}
		}
	}
}
