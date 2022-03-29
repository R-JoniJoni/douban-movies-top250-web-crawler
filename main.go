package main

import (
	"douban-movies-top250-web-crawler/db"
	"douban-movies-top250-web-crawler/node"
	"douban-movies-top250-web-crawler/page"
	"douban-movies-top250-web-crawler/relation"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

func main() {
	rand.Seed(25)
	movieUrl := `https://movie.douban.com/subject/\d*?/`
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.108 Safari/537.36",
		"Mozilla / 5.0(Macintosh; U; Intel Mac OS X 10_6_8; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"Mozilla/5.0 (Windows; U; Windows NT 6.1; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0",
		"Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.0; Trident/4.0)",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.0)",
		"Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.6; rv, 2.0.1) Gecko/20100101 Firefox/4.0.1",
		"Mozilla/5.0 (Windows NT 6.1; rv, 2.0.1) Gecko/20100101 Firefox/4.0.1",
		"Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.8.131 Version/11.11",
		"Opera/9.80 (Windows NT 6.1; U; en) Presto/2.8.131 Version/11.11",
		"Opera/9.80 (Windows NT 6.1; U; en) Presto/2.8.131 Version/11.11",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Maxthon 2.0)",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; TencentTraveler 4.0)",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1)",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; The World)",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Trident/4.0; SE 2.X MetaSr 1.0; SE 2.X MetaSr 1.0; .NET CLR 2.0.50727; SE 2.X MetaSr 1.0)",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; 360SE)",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Avant Browser)",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1)",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 UBrowser/6.2.4094.1 Safari/537.36",
	}

	// 爬取排名页面，一共10个，每爬到一个排名页面就搜索页面内的电影，并对其进行相应参数的提取，存到csv文件里
	for i := 0; i < 1; i++ {
		rankUrl := "https://movie.douban.com/top250" + "?start=" + strconv.Itoa(i*25) + "&filter="
		r := &page.Robot{
			Url:       rankUrl,
			UserAgent: userAgents[rand.Intn(len(userAgents))],
			SleepTime: 5 * time.Second,
		}
		body := r.Crawl()
		//fmt.Printf("%s\n", body[:16])

		// 使用正则表达式找到具体的电影页面
		regExp := regexp.MustCompile(movieUrl)
		urls := regExp.FindAllStringSubmatch(string(body), -1)
		fmt.Println("urls =", urls)

		crawlMoviesAndSaveTxt(urls, i, userAgents) // 爬具体的电影内容
	}

	directors := make(map[string]int, 0)
	actors := make(map[string]int, 0)
	films := make(map[string]int, 0)
	movieTypes := make(map[string]int, 0)

	node.GetNodes(directors, actors, films, movieTypes)         // 在txt中找到电影名、导演名、演员名、电影类型，并存为csv文件
	relation.GetRelations(directors, actors, films, movieTypes) // 在txt中找到导演电影关系、演员电影关系、合作关系、电影类型从属关系，并存为csv文件
	db.Import2Db()                                              // 把csv文件都导入到neo4j的import文件夹下
}

func crawlMoviesAndSaveTxt(urls [][]string, i int, userAgents []string) {
	for j := 0; j < 5; j++ {
		fmt.Println("2 * i =", i, "urls[2*i] =", urls[2*i], "urls[2*i][0] =", urls[2*i][0])
		movieRobot := &page.Robot{
			Url:       urls[2*j][0],
			SleepTime: 10 * time.Second,
			UserAgent: userAgents[rand.Intn(len(userAgents))],
		}
		movieBody := movieRobot.Crawl()
		movieBody = cutNewLine(movieBody)                                                           // 为了之后正则表达式匹配方便，去掉所有的换行符
		err := ioutil.WriteFile("data/contents/movie"+strconv.Itoa(i*25+j)+".txt", movieBody, 0644) // 把爬到的页面HTML存入txt文件中
		if err != nil {
			log.Fatalf("cannot write in file: %v", err)
		}
	}
}

func cutNewLine(body []byte) []byte {
	count := 0
	for _, v := range body {
		if v == '\n' {
			count++
		}
	}

	nb := make([]byte, len(body)-count)
	for i, j := 0, 0; i < len(nb); {
		if body[j] != '\n' {
			nb[i] = body[j]
			i++
		}
		j++
	}

	return nb
}
