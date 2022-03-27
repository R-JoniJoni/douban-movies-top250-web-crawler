package page

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Robot 类内包含有关爬取页面的功能
type Robot struct {
	Proxy     string
	SleepTime time.Duration
	Url       string
	UserAgent string
}

// Crawl 方法爬取页面的HTML，并返回
func (r *Robot) Crawl() []byte {
	req, err := http.NewRequest("GET", r.Url, nil) // 依据http包的规则，创建一个Request，之后利用它进行GET操作
	if err != nil {
		log.Fatalf("cannot new a request: %v\n", err)
	}

	req.Header.Add("User-Agent", r.UserAgent) // 处理反爬虫的手段1：设置HTTP报文中的User-Agent字段

	fmt.Println("Now crawling " + r.Url)
	client := &http.Client{}
	resp, err := client.Do(req) // 开始爬top250中的一页
	if err != nil {
		log.Fatalf("cannot do a request: %v\n", resp)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("cannot close body: %v\n", err)
		}
	}(resp.Body)

	time.Sleep(r.SleepTime) // 处理反爬虫的手段2：不要以过高的频率爬取

	body, err := ioutil.ReadAll(resp.Body) // 得到response报文中的body部分
	if err != nil {
		log.Fatalf("cannot read from the response's body: %v\n", err)
	}

	if len(body) == 0 {
		fmt.Printf("Failed to get url: %s, skip this and continue to do something else.\n", r.Url)
	}
	return body
}
