package utils

import (
	"database/sql"
	"fmt"
	"os/user"
	"strconv"
	"sync"

	"github.com/gocolly/colly"

	_ "github.com/mattn/go-sqlite3"
)

var (
	baseURL                  = "https://w12.mangafreak.net/"
	libraryURL               = "https://w12.mangafreak.net/Mangalist/All/"
	dirPath                  = "/Documents/mangafreak.db"
	osPath                   string
	maxConcurrentConnections = 16
)

func updateMangaListHandler(page int) {
	user, _ := user.Current()
	database, _ := sql.Open("sqlite3", user.HomeDir+"/"+dirPath)

	insert, _ := database.Prepare("INSERT INTO Manga(title, url) VALUES(?, ?)")

	c := colly.NewCollector()

	c.OnHTML(".list_item > .list_item_info > h3 > a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		insert.Exec(e.Text, link)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Printf("Finished Scraping Page %d\n", page)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "\nError:", err)
	})

	c.Visit(libraryURL + strconv.Itoa(page))
}

// This Function iterates through all the pages of mangafreak website and lists all the manga and updates mangafreak.db file
func UpdateMangaList() {
	user, _ := user.Current()
	database, _ := sql.Open("sqlite3", user.HomeDir+"/"+dirPath)

	createTable, _ := database.Prepare("CREATE TABLE IF NOT EXISTS Manga(id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT, url TEXT)")
	dropTable, _ := database.Prepare("DROP TABLE IF EXISTS Manga")

	dropTable.Exec()
	createTable.Exec()

	concManager := make(chan struct{}, maxConcurrentConnections)
	var wg sync.WaitGroup

	for page := 1; page <= 309; page++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			concManager <- struct{}{}
			updateMangaListHandler(page)
			<-concManager
		}(page)
	}

	wg.Wait()
}

func FetchMangaPageLink(id int) (string, string, error) {
	user, _ := user.Current()
	database, _ := sql.Open("sqlite3", user.HomeDir+"/"+dirPath)

	fetchRow, _ := database.Prepare("SELECT title, url FROM Manga WHERE id = ?")
	row := fetchRow.QueryRow(id)

	var (
		title string
		url   string
	)

	row.Scan(&title, &url)
	return title, url, nil
}

func SyncManga(title string, url string) {

	c := colly.NewCollector()

	concManager := make(chan struct{}, maxConcurrentConnections)
	var wg sync.WaitGroup

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	c.OnHTML(".manga_series_list > table > tbody > tr > td:last-child > a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		wg.Add(1)
		go func(title string, url string) {
			defer wg.Done()
			concManager <- struct{}{}
			DownloadFileHandler(title, url)
			<-concManager
		}(title, baseURL+link)

		wg.Wait()
	})

	c.Visit(baseURL + url)
}

func FetchFromDatabase(name string) {
	user, _ := user.Current()
	database, dbErr := sql.Open("sqlite3", user.HomeDir+"/"+dirPath)
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

	for rows.Next() {
		var (
			id    int
			title string
			url   string
		)
		rows.Scan(&id, &title, &url)
		fmt.Printf("%d : %s\n", id, title)
	}
	fmt.Println("Records Ended.")
}
