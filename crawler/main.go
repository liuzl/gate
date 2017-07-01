package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/crawlerclub/x/downloader"
	"github.com/crawlerclub/x/ds"
	"github.com/crawlerclub/x/parser"
	"github.com/crawlerclub/x/types"
	"github.com/golang/glog"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

var queue *ds.Queue

func initQueue() {
	var err error
	queue, err = ds.OpenQueue("./tasks")
	if err != nil {
		glog.Fatal(err)
	}
	defer queue.Close()

	file, err := os.Open("500w")
	if err != nil {
		glog.Fatal(err)
		return
	}
	defer file.Close()

	br := bufio.NewReader(file)
	for {
		line, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		queue.Enqueue(line)
	}
}

func work(proxy string, parseConf *types.ParseConf, wg *sync.WaitGroup, exitCh chan int) {
	defer wg.Done()
	glog.Info("start worker: ", proxy)
	defer glog.Info("exit worker: ", proxy)
	if queue == nil {
		glog.Error("queue is nil")
		return
	}
	fileName := "results/" + proxy + ".txt"
	fileName = strings.Replace(fileName, "http://", "", -1)
	fileName = strings.Replace(fileName, "https://", "", -1)
	f, err := os.Create(fileName)
	if err != nil {
		glog.Error("create file error:", fileName)
		return
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	for {
		select {
		case <-exitCh:
			return
		default:
			glog.Info(proxy, " work on next task")
			item, err := queue.Dequeue()
			if err == ds.ErrEmpty {
				glog.Error(err)
				time.Sleep(5 * time.Second)
				continue
			} else if err != nil {
				glog.Error(err)
				return
			}
			phone := string(item.Value)
			pageUrl := fmt.Sprintf("https://www.facebook.com/search/top/?q=%s", url.QueryEscape(phone))
			glog.Info(phone)
			req := &types.HttpRequest{Url: pageUrl, Method: "GET", Platform: "mobile", UseProxy: true, Proxy: proxy, Timeout: 5}
			res := downloader.Download(req)
			if res.Error != nil {
				glog.Error(res.Error)
				return
			}
			_, retItems, err := parser.Parse(res.Text, pageUrl, parseConf)
			ioutil.WriteFile("./html/"+phone+".html", res.Content, 0666)
			succ := false
			if len(retItems) > 0 {
				item := retItems[0]["results"].(map[string]interface{})
				if name, ok := item["name"]; ok {
					fmt.Println(name)
					succ = true
				}
			}

			str, _ := json.Marshal(retItems)
			fmt.Println(string(str))
			fmt.Fprintln(w, string(str))
			w.Flush()
			if !succ {
				return
			}
		}
	}
}

func stop(sigs chan os.Signal, exitCh chan int) {
	<-sigs
	glog.Info("receive stop signal")
	close(exitCh)
}

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
	confFb := flag.String("fb", "./fb.json", "fb configure file")
	mode := flag.String("mode", "run|queue", "mode")
	flag.Parse()
	if *mode == "queue" {
		initQueue()
		return
	}
	conf, _ := ioutil.ReadFile(*confFile)
	var parseConf types.ParseConf
	err := json.Unmarshal(conf, &parseConf)
	if err != nil {
		glog.Error(err)
		return
	}

	conf, _ = ioutil.ReadFile(*confFb)
	var fbParseConf types.ParseConf
	err = json.Unmarshal(conf, &fbParseConf)
	if err != nil {
		glog.Error(err)
		return
	}

	pageUrl := "https://free-proxy-list.net/"
	req := &types.HttpRequest{Url: pageUrl, Method: "GET", Platform: "pc", UseProxy: false}
	res := downloader.Download(req)
	if res.Error != nil {
		glog.Error(res.Error)
		return
	}
	_, retItems, err := parser.Parse(res.Text, pageUrl, &parseConf)
	f, err := os.Create(*outFile)
	if err != nil {
		log.Println("create file error:", *outFile)
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	if len(retItems) > 0 {
		queue, err = ds.OpenQueue("./tasks")
		if err != nil {
			glog.Fatal(err)
			return
		}
		defer queue.Close()

		exitCh := make(chan int)
		sigs := make(chan os.Signal)
		var wg sync.WaitGroup
		go stop(sigs, exitCh)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		i := 0
		proxies := retItems[0]["proxies"]
		for _, p := range proxies.([]interface{}) {
			proxy := p.(map[string]interface{})
			if ip, ok := proxy["Ip"]; ok {
				schema := "http://"
				if proxy["Https"].(string) == "yes" {
					schema = "https://"
				}
				one := fmt.Sprintf("%s%s:%s", schema, ip, proxy["Port"])
				valid := check(one)

				line := fmt.Sprintf("%s\t%s", one, valid)
				fmt.Println(line)
				fmt.Fprintln(w, line)
				w.Flush()
				if valid == "SUCC" {
					i += 1
					if i > 300 {
						break
					}
					wg.Add(1)
					go work(one, &fbParseConf, &wg, exitCh)
				}
			}
		}
		wg.Wait()
	}
}
