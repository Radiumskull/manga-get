<h3> Manga-Get </h3>

This is a web-scrapper built using Golang which scrapes MangaFreak to maintain a database. It can then download all the manga chapters of a selected Manga at once.

Now it contains 5000 Manga which it can download.

Third-Party Packages Used:
<ul>
<li> github.com/gocolly/colly
<li> github.com/mattn/go-sqlite3
</ul>

Future Updates:
1. I wish to package this into CLI app in future.
2. Sometimes the Mangafreak url changes which needs to be dynamic.
3. The Scrapping to not Concurrent. It takes around 309secs to update database.

Developer: Aritra Bhattacharjee(RadiumSkull)
