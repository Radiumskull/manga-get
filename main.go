package main

import (
    "fmt"
    "manget/utils"
)


var (
    baseURL = "https://w12.mangafreak.net/"
    libraryURL = "https://w12.mangafreak.net/Mangalist/All/"
)

func main(){

    //utils.UpdateMangaList()
    //utils.FetchFromDatabase("Akame")

    title, url, err := utils.FetchMangaPageLink(3345)
    if err != nil {
        fmt.Println(err.Error())
    }

    utils.CreateDir(title)
    utils.SyncManga(title, baseURL + url)
}
