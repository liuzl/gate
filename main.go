package main

import (
	"flag"
	"fmt"
	"github.com/crawlerclub/x/downloader"
	"github.com/crawlerclub/x/types"
	"log"
)

func main() {
	proxy := flag.String("proxy", "http://127.0.0.1:1080", "proxy address")
	flag.Parse()
	fmt.Println("crawl free proxies from the web")
	pageUrl := "https://free-proxy-list.net/"
	req := &types.HttpRequest{Url: pageUrl, Method: "GET", Platform: "pc", UseProxy: false, Proxy: *proxy}
	fmt.Println(pageUrl)
	res := downloader.Download(req)
	if res.Error != nil {
		log.Println(res.Error)
	}
	log.Println(res.Text)
}
