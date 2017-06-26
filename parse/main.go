package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/crawlerclub/x/parser"
	"github.com/crawlerclub/x/types"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	confFile := flag.String("conf", "./fb.json", "crawler configure file")
	flag.Parse()
	conf, _ := ioutil.ReadFile(*confFile)
	var parseConf types.ParseConf
	err := json.Unmarshal(conf, &parseConf)
	if err != nil {
		log.Fatal(err)
	}

	pageUrl := "https://www.facebook.com/search/top/?q=18232032959"
	Text, _ := ioutil.ReadFile("fb.html")
	html := strings.Replace(string(Text), "<head>", "<head><meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\">", 1)
	_, retItems, err := parser.Parse([]byte(html), pageUrl, &parseConf)
	//str, _ := json.Marshal(retItems)
	fmt.Println(retItems[0]["results"])
	//fmt.Println(string(str))
}
