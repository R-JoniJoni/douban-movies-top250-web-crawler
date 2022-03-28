package relation

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

func GetRelations(
	directors map[string]int,
	actors map[string]int,
	films map[string]int,
	movieTypes map[string]int,
) {
	fmt.Println("Now getting Relations.")

	actedIn := make(map[int][]int, 0)     // film -> []actor
	belongTo := make(map[int][]int, 0)    // film -> []type
	cooperation := make(map[int][]int, 0) // actor -> []director
	directed := make(map[int][]int, 0)    // film -> []director

	for i := 0; i < 250; i++ {
		// 读取moviei.txt
		file, err := os.Open("data/contents/movie" + strconv.Itoa(i) + ".txt")
		if err != nil {
			fmt.Println("Failed to open movie" + strconv.Itoa(i) + ".txt")
			continue
		}
		moviePage, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatalf("cannot read all from file: %v\n", err)
		}

		para := GetParametersInPage(moviePage, films, directors, actors, movieTypes) // 找到电影名、导演名、演员名、类型名(都是对应的id)
		saveInMap(para, actedIn, 0, 2)                                               // 把film-actor的关系写进map中
		saveInMap(para, belongTo, 0, 3)                                              // 把film-type的关系写进map中
		saveInMap(para, cooperation, 2, 1)                                           // 把actor-director的关系写进map中
		saveInMap(para, directed, 0, 1)                                              // 把film-director的关系写进map中
	}

	fmt.Println("Finished getting Relations. Now saving Relations to csv files.")

	saveRelation2Csv(actedIn, "acted_in", "出演")
	saveRelation2Csv(belongTo, "belong_to", "属于")
	saveRelation2Csv(cooperation, "cooperation", "合作")
	saveRelation2Csv(directed, "directed", "拍了")

	fmt.Println("Finishen saving Relations to files.")
}

func saveRelation2Csv(myMap map[int][]int, fileName string, relation string) {
	csvFile, err := os.Create("data/details/" + fileName + ".csv")
	if err != nil {
		log.Fatalf("cannot open csv file %s: %v\n", fileName, err)
	}
	defer csvFile.Close()

	// 获取csv的Writer
	writer := csv.NewWriter(csvFile)
	// 写入csv文件
	for k, vs := range myMap {
		for _, v := range vs {
			oneLine := []string{strconv.Itoa(k), strconv.Itoa(v), relation}
			err := writer.Write(oneLine)
			if err != nil {
				log.Fatalf("cannot write relations to csv file: %v\n", err)
			}
		}
	}

	// 确保所有内存数据刷到csv文件
	writer.Flush()
}

func saveInMap(para [4][]int, actedIn map[int][]int, kIndex int, vIndex int) {
	for _, v := range para[kIndex] {
		if _, ok := actedIn[v]; !ok {
			actedIn[v] = para[vIndex]
		} else {
			actedIn[v] = append(actedIn[v], para[vIndex]...)
		}
	}
}

// GetParametersInPage 函数返回输入page中的4种参数，依次为：电影名、导演、演员、类型(均为其对应的id)
func GetParametersInPage(
	page []byte,
	films map[string]int,
	directors map[string]int,
	actors map[string]int,
	types map[string]int,
) (para [4][]int) {
	// 找电影名
	regExp := regexp.MustCompile(`<title> {8}.*? \(豆瓣\)</title>`)
	film := regExp.Find(page)
	filmStr := string(film)
	filmStr = strings.TrimPrefix(filmStr, `<title>        `)
	filmStr = strings.TrimSuffix(filmStr, ` (豆瓣)</title>`)
	para[0] = []int{films[filmStr]}

	// 找导演名
	regExp = regexp.MustCompile(`"director":.*?]`)
	directorNoise := regExp.Find(page)
	regExp = regexp.MustCompile(`"name": ".*?"`)
	director := regExp.FindAll(directorNoise, -1)
	var tmp []int
	for _, v := range director {
		directorStr := string(v)
		directorStr = strings.TrimPrefix(directorStr, `"name": "`)
		directorStr = strings.TrimSuffix(directorStr, `"`)
		tmp = append(tmp, directors[directorStr])
	}
	para[1] = tmp

	// 找演员名
	regExp = regexp.MustCompile(`"actor":.*?]`)
	actorNoise := regExp.Find(page)
	regExp = regexp.MustCompile(`"name": ".*?"`)
	actor := regExp.FindAll(actorNoise, -1)
	tmp = make([]int, 0)
	for _, v := range actor {
		actorStr := string(v)
		actorStr = strings.TrimPrefix(actorStr, `"name": "`)
		actorStr = strings.TrimSuffix(actorStr, `"`)
		tmp = append(tmp, actors[actorStr])
	}
	para[2] = tmp

	// 找类型名
	regExp = regexp.MustCompile(`<span property="v:genre">.*?</span>`)
	movieType := regExp.FindAll(page, -1)
	tmp = make([]int, 0)
	for _, v := range movieType {
		typeStr := string(v)
		typeStr = strings.TrimPrefix(typeStr, `<span property="v:genre">`)
		typeStr = strings.TrimSuffix(typeStr, `</span>`)
		tmp = append(tmp, types[typeStr])
	}
	para[3] = tmp

	return
}
