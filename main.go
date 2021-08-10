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

// func main(){
//     fmt.Printf("\n\n -------------------------------------------------------\n")
//     fmt.Printf("1. Search Manga\n2. Download Manga\n3. Update Manga Database\n4. Exit\n*If its your first time, Update Manga Database\n")
//     fmt.Printf("\nEnter Choice: ")

//     for {
//         var choice string;

//         fmt.Scanln(&choice)

//         switch choice {
//         case "1":
//             var input string
//             fmt.Printf("Enter Manga Name: ")
//             fmt.Scanf("%s", &input)

//             utils.FetchFromDatabase(input)
//             break;
//         case "2":
//             var input string
//             fmt.Printf("Enter Manga ID: ")
//             fmt.Scanf("%s", &input)

//             mangaId, _ := strconv.Atoi(input)
// title, url, err := utils.FetchMangaPageLink(mangaId)
// if err != nil{
//     fmt.Println(err.Error())
//     os.Exit(3)
// }

//             utils.CreateDir("./Manga")
//             utils.CreateDir("./Manga/"+ title)
//             utils.SyncManga(title, url)

//             fmt.Println("Enter 4 to Exit.")
//             break;
//         case "3":
//             utils.UpdateMangaList()
//             break;
//         case "4":
//             fmt.Println("Enjoy Reading Manga. :)")
//             os.Exit(3)
//             break;
//         default:
//             fmt.Printf("\nEnter Choice: ")
//             break;
//         }
//     }
// }
