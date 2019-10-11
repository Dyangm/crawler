package command

import (
	"fmt"
	"github.com/Dyangm/crawler/config"
	"github.com/Dyangm/crawler/downloader"
	"github.com/Dyangm/crawler/fetch"
	"github.com/Dyangm/crawler/search"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Handler struct {
	urlArr     []*Request
	kCountHttp int

	command    string
}

func NewHandler() *Handler {
	h := &Handler{}
	return h
}

func (h *Handler) CommandHandler() {
	for {
		var cmd string = "0"
		if len(h.urlArr) > 0 {
			fmt.Println("0:搜索;1:搜索信息;2:选择下载;否则退出")
			fmt.Scan(&cmd)
		}

		switch cmd {
		case "0":
			h.searchCommand()
		case "1":
			h.printAllData()
		case "2":
			h.downloaderCommand(0)
		default:
			os.Exit(0)
		}
	}
}

type Request struct {
	City string
	Url  string
}

func (h *Handler) searchCommand() {
	fmt.Println("请输入搜索关键词:")
	var name string
	fmt.Scan(&name)
	searchHandler := &search.SearchHandler{}
	config, _ := config.GetConfig()
	for _, v := range config.NovelWebInfo {
		searchHandler.Search(name, v)
	}

	h.urlArr = make([]*Request, 0)
	for k, v := range searchHandler.SearchMap {
		request := &Request{}
		request.Url = v
		request.City = k
		h.urlArr = append(h.urlArr, request)
	}
	fmt.Printf("搜索结果: %d 条数据\n", len(h.urlArr))
	if len(h.urlArr) == 0 {
		return
	}
	h.printAllData()
}

const kPrintNumber = 10

func (h *Handler) printAllData() {
	var command int
	if len(h.urlArr) == 0 {
		fmt.Println("查询数据不存在，请从新输入查询！")
		return
	}
	fmt.Println("查询全部数据;")
	num := kPrintNumber
	for {
		count := h.kCountHttp + num
		if count > len(h.urlArr) {
			count = len(h.urlArr)
		} else if h.kCountHttp > count {
			h.kCountHttp = count
		} else if h.kCountHttp < 0 {
			h.kCountHttp = 0
		}
		for i := h.kCountHttp; i < count; i++ {
			fmt.Printf("%d： %s\n", i, h.urlArr[i].City)
		}
		fmt.Printf("1:下载; 2: 选择查看第N条数据;")
		if len(h.urlArr) > count {
			h.kCountHttp = count
			fmt.Printf("3: 下一页; ")
		}
		if h.kCountHttp > num {
			fmt.Printf("0: 上一页; ")
		}
		fmt.Println(" 否则返回;")
		fmt.Scan(&command)

		switch command {
		case 1:
			h.downloaderCommand(command)
		case 2:
			h.printData(command)
		case 0:
			h.kCountHttp -= num * 2
		case 3:
			break
		default:
			return
		}
	}
}

func (h *Handler) printData(command int) {
	fmt.Println("查询第N条数据详情: ")
	fmt.Scan(&command)
	fmt.Println(h.urlArr[command].City)
	fetch, e := fetcher.FetchMethodGet(h.urlArr[command].Url)
	if e != nil {
		log.Error(e)
	}
	html := string(fetch)
	html = search.ConvertToString(html, "gbk", "utf8")
	document, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Error(err)
	}

	msgs :=document.Find("div>h3").Text()
	fmt.Println(msgs)
	dataArr := downloader.FindAllChaptersFromUrl(h.urlArr[command].Url, "div>ul.list-charts")
	fmt.Println("章节总数: ", len(dataArr))
	fmt.Print("作品简介:	")
	aNode :=document.Find("div>p>br")
	for _, v := range aNode.Nodes {
		if strings.Count(v.PrevSibling.Data, "&") > 0 {
			fmt.Println(v.PrevSibling.Data)
			continue
		}
		fmt.Println(v.PrevSibling.Data)
	}
	msg := aNode.Nodes[len(aNode.Nodes)-1].NextSibling.Data
	if msg != "" {
		index := strings.Index(msg, "(")
		fmt.Println(msg[:index])
	}
	fmt.Println()
}

func (h *Handler) downloaderCommand(command int) {
	fmt.Println("下载小说序列号: ")
	fmt.Scan(&command)
	name := h.urlArr[command].City
	dataArr := downloader.FindAllChaptersFromUrl(h.urlArr[command].Url, "div>ul.list-charts")
	downloader.DownloadAllValidChapters(name, dataArr)
}
