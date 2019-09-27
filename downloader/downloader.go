package downloader

import (
	"fmt"
	"github.com/Dyangm/crawler/fetch"
	"github.com/Dyangm/crawler/search"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type (
	Downloader struct {
	}
	Chapter struct {
		Name string
		Url  string
	}
)

func FindAllChaptersFromUrl(url, listReg string) (chapters []Chapter) {
	bytes, err := fetcher.Fetch(url)
	if err != nil {
		log.Println("Fatal error: ", err.Error())
	}

	html := string(bytes)
	html = search.ConvertToString(html, "gbk", "utf8")

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Println("Fatal error: ", err.Error())
	}
	// "div.listmain"
	doc.Find(listReg).Each(func(index int, s *goquery.Selection) {
		s.Find("a").Each(func(i int, s *goquery.Selection) {
			src, finded := s.Attr("href")
			name := s.Text()
			if finded {
				chapters = append(chapters, Chapter{
					name,
					"http://www.shu05.com" + src,
				})
			}
		})
	})
	return
}

func DownloadAllValidChapters(name string, chapters []Chapter) (lastDownload, lastName string) {
	for _, chap := range chapters {
		fmt.Printf("下载：%s ......", chap.Name)
		content := downloadChapterContent(chap.Url)
		if content == "" {
			return
		}
		saveToFile(chap.Name, content, name+".txt")
		lastDownload = chap.Url
		lastName = chap.Name
		fmt.Println("完成。")
	}
	return
}

func downloadChapterContent(chapterUrl string) (content string) {
	bytes, err := fetcher.Fetch(chapterUrl)
	if err != nil {
		log.Println("Fatal error: ", err.Error())
	}

	html := string(bytes)
	html = search.ConvertToString(html, "gbk", "utf8")
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Println("error:", err.Error())
	}

	doc.Find("body>div").Each(func(i int, s *goquery.Selection) {
		aNode := s.Find("br")
		for _, v := range aNode.Nodes {
			if strings.Count(v.PrevSibling.Data, "&") > 0 {
				fmt.Println(v.PrevSibling.Data)
				continue
			}
			content += v.PrevSibling.Data
		}
	})

	return
}

func filterValidChapters(chapters []Chapter, lastNum int) []Chapter {
	if lastNum == 0 {
		return chapters
	}

	return chapters[lastNum+1:]
}

func saveToFile(title, content, filePath string) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		defer file.Close()
		file.WriteString(fmt.Sprintf("\n%s\n", title))
		file.WriteString(content)
	} else {
		fmt.Println("打开文件错误:", err.Error())
	}
}
