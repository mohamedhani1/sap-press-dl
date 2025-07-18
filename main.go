package main

import (
	"fmt"
	"log"
	"bufio"
	"os"

	"sappress/sappress"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "SapPress-DL",
		Usage: "Download True EPUB Books\nCreated by t.me/@Caliginous_0",
		Commands: []*cli.Command{
			{
				Name:  "download",
				Usage: "Download a new book by ID",
				Action: func(c *cli.Context) error {
					var bookID string
					bookID = c.String("bookid")
					for {
						
						if bookID == "" {
							fmt.Printf("Please add the BookID (xxxx): ")
							fmt.Scan(&bookID)
						}
						
						sappress.CheckToken()

						// Load config and setup downloader
						config, err := sappress.LoadConfig()
						if err != nil {
							return fmt.Errorf("failed to load config: %w", err)
						}

						downloader := &sappress.Downloader{
							Config:     config,
							HttpClient: sappress.NewAuthenticatedClient(config.Token),
							Threads:    c.Int("threads"),
						}

						downloader.Download(bookID)

						bookID = ""
					}
				
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "bookid",
						Usage:    "ID of the book to download",
						Required: false,
					},
					&cli.IntFlag{
						Name:  "threads",
						Usage: "Number of concurrent threads",
						Value: 16,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
	fmt.Print("Press Enter to exit...")
    bufio.NewReader(os.Stdin).ReadBytes('\n')
}
