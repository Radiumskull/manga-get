package utils

import (
    "fmt"
    "github.com/gocolly/colly"
    "database/sql"
    "strconv"

    _ "github.com/mattn/go-sqlite3"
)

var (
    baseURL = "https://w12.mangafreak.net/"
    libraryURL = "https://w12.mangafreak.net/Mangalist/All/"
)

func updateMangaListHandler(page int){
    database, _ := sql.Open("sqlite3", "./data/mangafreak.db")
    insert, _ := database.Prepare("INSERT INTO Manga(title, url) VALUES(?, ?)")

    c := colly.NewCollector()


    c.OnHTML(".list_item > .list_item_info > h3 > a[href]", func(e *colly.HTMLElement) {
        link := e.Attr("href")
        //fmt.Println(reflect.TypeOf(e.Text), reflect.TypeOf(link))
        insert.Exec(e.Text, link);
    })

    c.OnScraped(func (r *colly.Response){
        fmt.Printf("Finished Scraping Page %d\n", page)
        //wg.Done()
    })

    c.OnError(func(r *colly.Response, err error){
        fmt.Println("Request URL:", r.Request.URL, "\nError:", err)
    })


    c.Visit(libraryURL + strconv.Itoa(page))
}


// This Function iterates through all the pages of mangafreak website and lists all the manga and updates mangafreak.db file
func UpdateMangaList(){
    //var wg sync.WaitGroup
    database, _ := sql.Open("sqlite3", "./data/mangafreak.db")

    createTable, _ := database.Prepare("CREATE TABLE IF NOT EXISTS Manga(id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT, url TEXT)");
    dropTable, _ := database.Prepare("DROP TABLE IF EXISTS Manga");

    dropTable.Exec();
    createTable.Exec();
    for page := 1; page <= 309; page++ {
        //wg.Add(1)
        updateMangaListHandler(page);
    }

    //wg.Wait()
}

func FetchMangaPageLink(id int) (string, string, error){
    database, _ := sql.Open("sqlite3", "./data/mangafreak.db")

    fetchRow, _ := database.Prepare("SELECT title, url FROM Manga WHERE id = ?")
    row := fetchRow.QueryRow(id)

    var (
        title string
        url string
    )

    row.Scan(&title, &url)
    return title, url, nil

}

func SyncManga(title string, url string) {

    c := colly.NewCollector()

    c.OnRequest(func (r *colly.Request) {
        fmt.Println("Visiting", r.URL.String())
    })

    c.OnScraped(func (r *colly.Response){
        fmt.Println("Finished", r.Request.URL)
    })

    c.OnHTML(".manga_series_list > table > tbody > tr > td:last-child > a", func(e *colly.HTMLElement) {
        link := e.Attr("href")

        DownloadFileHandler(title, baseURL + link)
    })

    c.Visit(url)
}

func FetchFromDatabase(name string){
    database, dbErr := sql.Open("sqlite3", "./data/mangafreak.db")
    if dbErr != nil {
        fmt.Println("Database Not Found.")
        panic(dbErr)
    }
    searchTable, _ := database.Prepare("SELECT id, title, url FROM Manga WHERE title LIKE ?")
    rows, queryErr := searchTable.Query("%" + name + "%")
    if queryErr != nil {
        fmt.Println("Query Err")
        panic(queryErr)
    }

    for rows.Next(){
        var (
            id int
            title string
            url string
        )
        rows.Scan(&id, &title, &url)
        fmt.Printf("%d : %s\n", id, title)
    }
    fmt.Println("Records Ended.")
}

