package node

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func GetNodes(directors map[string]int, actors map[string]int, films map[string]int, movieTypes map[string]int) {
	fmt.Println("Now getting Nodes.")

	id := 0

	for i := 0; i < 250; i++ { // 循环在每个txt文件中找
		file, err := os.Open("data/contents/movie" + strconv.Itoa(i) + ".txt")
		if err != nil {
			fmt.Println("Failed to open movie" + strconv.Itoa(i) + ".txt")
			continue
		}
		moviePage, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatalf("cannot read all from file: %v\n", err)
		}

		// 找电影名
		regExp := regexp.MustCompile(`<title> {8}.*? \(豆瓣\)</title>`)
		film := regExp.Find(moviePage)
		filmStr := string(film)
		filmStr = strings.TrimPrefix(filmStr, `<title>        `)
		filmStr = strings.TrimSuffix(filmStr, ` (豆瓣)</title>`)
		if _, ok := films[filmStr]; !ok {
			films[filmStr] = id
			id++
		}

		// 找导演名
		regExp = regexp.MustCompile(`"director":.*?]`)
		directorNoise := regExp.Find(moviePage)
		regExp = regexp.MustCompile(`"name": ".*?"`)
		director := regExp.FindAll(directorNoise, -1)
		for _, v := range director {
			directorStr := string(v)
			directorStr = strings.TrimPrefix(directorStr, `"name": "`)
			directorStr = strings.TrimSuffix(directorStr, `"`)
			if _, ok := directors[directorStr]; !ok {
				directors[directorStr] = id
				id++
			}
		}

		// 找演员名
		regExp = regexp.MustCompile(`"actor":.*?]`)
		actorNoise := regExp.Find(moviePage)
		regExp = regexp.MustCompile(`"name": ".*?"`)
		actor := regExp.FindAll(actorNoise, -1)
		for _, v := range actor {
			actorStr := string(v)
			actorStr = strings.TrimPrefix(actorStr, `"name": "`)
			actorStr = strings.TrimSuffix(actorStr, `"`)
			if _, ok := actors[actorStr]; !ok {
				actors[actorStr] = id
				id++
			}
		}

		// 找类型名
		regExp = regexp.MustCompile(`<span property="v:genre">.*?</span>`)
		movieType := regExp.FindAll(moviePage, -1)
		for _, v := range movieType {
			typeStr := string(v)
			typeStr = strings.TrimPrefix(typeStr, `<span property="v:genre">`)
			typeStr = strings.TrimSuffix(typeStr, `</span>`)
			if _, ok := movieTypes[typeStr]; !ok {
				movieTypes[typeStr] = id
				id++
			}
		}

		file.Close()
	}

	// 存csv文件
	saveNodes2Csv(films, "film")
	saveNodes2Csv(actors, "actor")
	saveNodes2Csv(directors, "director")
	saveNodes2Csv(movieTypes, "type")

	fmt.Println("Finished getting Nodes.")
}

func saveNodes2Csv(myMap map[string]int, fileName string) {
	csvFile, err := os.Create("data/details/" + fileName + ".csv")
	if err != nil {
		log.Fatalf("cannot open csv file %s: %v\n", fileName, err)
	}
	defer csvFile.Close()

	// 获取csv的Writer
	writer := csv.NewWriter(csvFile)
	// 并写入csv文件
	for name, id := range myMap {
		oneLine := []string{name, strconv.Itoa(id)}
		err := writer.Write(oneLine)
		if err != nil {
			log.Fatalf("cannot write csv file: %v\n", err)
		}
	}

	// 确保所有内存数据刷到csv文件
	writer.Flush()
}
