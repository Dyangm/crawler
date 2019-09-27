package command

import (
	"fmt"
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
		var cmd string
		fmt.Println("1:搜索;2:搜索信息;否则退出")
		fmt.Scan(&cmd)
		switch cmd {
		case "1":
			h.searchCommand()
		case "2":
			h.printAllData()
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
	searchHandler.Search(name)
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
			fmt.Printf("%d %s\n", i, h.urlArr[i].City)
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
			go h.downloaderCommand(command)
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
	log.Println("查询第N条数据详情: ")
	fmt.Scan(&command)
	log.Println(h.urlArr[command].City, h.urlArr[command].Url)
	fetch, e := fetcher.Fetch(h.urlArr[command].Url)
	if e != nil {
		log.Error(e)
	}
	html := string(fetch)
	html = search.ConvertToString(html, "gbk", "utf8")
	document, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Error(err)
	}

	content := strings.Replace(document.Find("").Text(), "", "", -1)
	fmt.Println(content)
}

func (h *Handler) downloaderCommand(command int) {
	fmt.Println("下载小说序列号: ")
	fmt.Scan(&command)
	name := h.urlArr[command].City
	dataArr := downloader.FindAllChaptersFromUrl(h.urlArr[command].Url, "div>ul.list-charts")
	downloader.DownloadAllValidChapters(name, dataArr)
}
