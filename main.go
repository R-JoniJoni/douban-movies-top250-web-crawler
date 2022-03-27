package main

import (
	"douban-movies-top250-web-crawler/page"
	"strconv"
	"time"
)

func main() {

	// 爬取排名页面，一共10个
	for i := 0; i < 10; i++ {
		rankUrl := "https://movie.douban.com/top250" + "?start=" + strconv.Itoa(i*25) + "&filter="
		r := &page.Robot{
			Url:       rankUrl,
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.108 Safari/537.36",
			SleepTime: 5 * time.Second,
		}
		body := r.Crawl()
		//fmt.Printf("%s\n", body[:16])

	}
}
