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

// CrawlRankingPage 方法爬取top250大页面的一页HTML，不是具体的某个属于top250的某一部电影对应页面的HTML，它返回页面的HTML
func (r *Robot) CrawlRankingPage() []byte {
	req, err := http.NewRequest("GET", r.Url, nil) // 依据http包的规则，创建一个Request，之后利用它进行GET操作
	if err != nil {
		log.Fatalf("cannot new a request: %v\n", err)
	}

	req.Header.Add("User-Agent", r.UserAgent) // 处理反爬虫的手段1：设置HTTP报文中的User-Agent字段

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

	return body
}

// CrawlMoviePage 方法爬取具体的某一部电影页面的HTML，并把其存入 ./data/contents/ 下的一个txt文件中
func (r *Robot) CrawlMoviePage() error {
	return fmt.Errorf("NOT A FINISHED METHOD")
}
