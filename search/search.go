package search

import (
	"bytes"
	"github.com/Dyangm/crawler/config"
	"github.com/Dyangm/crawler/fetch"
	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"strings"
)

type Request struct {
	Name string
	Url  string
}

type SearchHandler struct {
	SearchMap map[string]string
	ListMap map[string]string
}

func (hand *SearchHandler) Search(name string, webInfo config.WebInfo) error {
	val, _ := Utf8ToGbk([]byte(name))
	curl := "http://www.shu05.com/modules/article/search.php?searchkey=" + string(val)
	bytes, err := fetcher.FetchMethodPost(curl)
	if err != nil {
		return err
	}
	str := string(bytes)
	hand.SearchMap, err = hand.parser(str, webInfo.SearchReg)
	if err != nil {
		return err
	}

	if len(hand.SearchMap) == 0 {
		return nil
	}
	searchPage, err := hand.parser(str, webInfo.SearchPageReg)
	if err != nil {
		return err
	}

	for _, v := range searchPage {
		url := v
		if !strings.Contains(v, webInfo.Homepage) {
			url = webInfo.Homepage + v
		}
		fetch, err := fetcher.FetchMethodGet(url)
		str := string(fetch)
		searchMap, err := hand.parser(str, webInfo.SearchReg)
		if err != nil {
			return err
		}

		for k, v := range searchMap {
			hand.SearchMap[k] = v
		}
	}

	return nil
}

func (hand *SearchHandler) parser(str, selector string) (map[string]string, error) {
	str = ConvertToString(str, "gbk", "utf-8")
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(str))
	if err != nil {
		return nil, err
	}

	searchMap := make(map[string]string)
	dom.Find("body>div").Each(func(i int, s *goquery.Selection) {
		aNode := s.Find(selector)
		for _, v := range aNode.Nodes {
			if len(v.Attr) != 1 {
				continue
			}
			searchMap[v.FirstChild.Data] = v.Attr[0].Val
		}
	})

	return searchMap, nil
}

func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
