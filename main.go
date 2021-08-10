package main

import (
	"errors"
	"log"
	"manget/utils"
	"os"
	"strconv"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			&cli.Command{
				Name:  "search",
				Usage: "manget search Akame ga kill",
				Action: func(c *cli.Context) error {
					if c.NArg() != 1 {
						return errors.New("Invalid Syntax")
					}
					input := c.Args().Get(0)

					utils.FetchFromDatabase(input)
					return nil
				},
			},
			&cli.Command{
				Name:  "sync",
				Usage: "manget sync",
				Action: func(_ *cli.Context) error {
					utils.UpdateMangaList()
					return nil
				},
			},
			&cli.Command{
				Name:  "download",
				Usage: "manget download 201",
				Action: func(c *cli.Context) error {
					if c.NArg() != 1 {
						return errors.New("Invalid Syntax")
					}

					mangaId, inpErr := strconv.Atoi(c.Args().Get(0))
					if inpErr != nil {
						return errors.New("Enter an valid Manga Id which is a number.")
					}
					title, url, err := utils.FetchMangaPageLink(mangaId)
					if err != nil {
						return errors.New("Manga Id not valid. Re-run manget search and enter valid id")
					}
					utils.CreateDir("./Manga")
					utils.CreateDir("./Manga/" + title)
					utils.SyncManga(title, url)

					return nil
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
