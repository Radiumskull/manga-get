package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func CreateDir(path string) {
	pwd, _ := os.Getwd()
	currDir := pwd + "/"

	_, err := os.Stat(currDir + path)
	if err == nil {
		return
	}
	mkErr := os.Mkdir(currDir+path, 0755)
	if mkErr != nil {
		check(mkErr)
	}
}

func FileNameParse(url string) (string, error) {
	components := strings.Split(url, "/")
	n := len(components)
	fileName := strings.Replace(components[n-1], "_", " ", -1)
	//fmt.Printf("%s\n", fileName)
	return fileName, nil
}

func DownloadFileHandler(dirName string, url string) {
	pwd, _ := os.Getwd()
	currDir := pwd + "/"

	resp, _ := http.Get(url)
	defer resp.Body.Close()

	fileName, err := FileNameParse(url)
	check(err)
	filePath := currDir + "Manga/" + dirName + "/" + fileName + ".cbr"
	out, _ := os.Create(filePath)

	defer out.Close()

	_, _ = io.Copy(out, resp.Body)
	fmt.Printf("%s Downloaded\n", fileName)
}
