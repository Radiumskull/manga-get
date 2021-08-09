package main

import (
    "fmt"
    "os"
    "manget/utils"
    "strconv"
)

func main(){
    for {
        var choice string;

        fmt.Printf("1. Search Manga\n2. Download Manga\n3. Update Manga Database\n4. Exit\n*If its your first time, Update Manga Database\n")
        fmt.Scanln(&choice)

        switch choice {
        case "1":
            var input string
            fmt.Printf("Enter Manga Name: ")
            fmt.Scanf("%s", &input)

            utils.FetchFromDatabase(input)
            break;
        case "2":
            var input string
            fmt.Printf("Enter Manga ID: ")
            fmt.Scanf("%s", &input)

            mangaId, _ := strconv.Atoi(input)
            title, url, err := utils.FetchMangaPageLink(mangaId)
            if err != nil{
                fmt.Println(err.Error())
                os.Exit(3)
            }

            utils.CreateDir("./Manga")
            utils.CreateDir("./Manga/"+ title)
            utils.SyncManga(title, url)
            break;
        case "3":
            utils.UpdateMangaList()
            break;
        case "4":
            os.Exit(3)
            break;
        default:
            os.Exit(3)
        }
    }
}
