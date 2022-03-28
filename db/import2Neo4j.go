package db

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func Import2Db() {
	clearFiles()      // 先把neo4j的import文件夹下的文件都删掉
	copyNeededFiles() // 把data/details文件夹下的文件复制到neo4j的import文件夹下
}

func clearFiles() {
	dir, err := ioutil.ReadDir("D:/Neo4j/neo4j-community-4.3.7-windows/neo4j-community-4.3.7/import")
	if err != nil {
		log.Fatalf("cannot read directory: %v\n", err)
	}

	for _, d := range dir {
		err := os.RemoveAll(path.Join([]string{"D:/Neo4j/neo4j-community-4.3.7-windows/neo4j-community-4.3.7/import", d.Name()}...))
		if err != nil {
			log.Fatalf("cannot remove: %v\n", err)
		}
	}
}

func copyNeededFiles() {
	srcPath := "data/details/"
	desPath := "D:/Neo4j/neo4j-community-4.3.7-windows/neo4j-community-4.3.7/import/"
	fileNames := [8]string{
		"acted_in.csv",
		"actor.csv",
		"belong_to.csv",
		"cooperation.csv",
		"directed.csv",
		"director.csv",
		"film.csv",
		"type.csv",
	}

	for i := 0; i < 8; i++ {
		_, err := copyFile(srcPath+fileNames[i], desPath+fileNames[i])
		if err != nil {
			log.Fatalf("cannot copy file "+fileNames[i]+": %v\n", err)
		}
	}
}

func copyFile(srcFile, destFile string) (int64, error) {
	file1, err := os.Open(srcFile)
	if err != nil {
		return 0, err
	}
	file2, err := os.OpenFile(destFile, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return 0, err
	}
	defer file1.Close()
	defer file2.Close()

	return io.Copy(file2, file1)
}
