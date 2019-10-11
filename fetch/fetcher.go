package fetcher

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"io"
	"io/ioutil"
	"net/http"
)

func FetchMethodGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	/*if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status code : %d", resp.StatusCode)
	}

	e := determineEncoding(resp.Body)
	utf8Reader := transform.NewReader(resp.Body, e.NewDecoder())*/

	return ioutil.ReadAll(resp.Body)
}

func FetchMethodPost(url string) ([]byte, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //这个一定要加，不加form的值post不过去，被坑了两小时
	req.Header.Add("Charset", "UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func determineEncoding(r io.Reader) encoding.Encoding {
	bytes, err := bufio.NewReader(r).Peek(1024)
	if err != nil {
		logrus.Error("Fetcher error : ", err)
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")

	return e
}
