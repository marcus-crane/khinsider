package khinsider

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"

	"github.com/marcus-crane/khinsider/pkg/indexer"
	"github.com/marcus-crane/khinsider/pkg/scrape"
	"github.com/marcus-crane/khinsider/pkg/search"
)

func Execute() {
	app := &cli.App{
		Name:     "khinsider",
		Version:  "2.0.0",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Marcus Crane",
				Email: "khinsider@utf9k.net",
			},
		},
		Usage: "khinsider - easily fetch videogame soundtracks from downloads.khinsider.com",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "debug", Aliases: []string{"d"}},
		},
		Before: func(c *cli.Context) error {
			if c.Bool("debug") {
				pterm.EnableDebugMessages()
			}
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "search",
				Aliases: []string{"s"},
				Usage:   "search for an album to download",
				Before: func(c *cli.Context) error {
					indexExists := indexer.CheckIndexExists()
					if indexExists {
						pterm.Debug.Println("Checking for updates")
						updateAvailable := indexer.CheckIndexUpdateAvailable()
						if updateAvailable {
							err := indexer.DownloadIndex()
							if err != nil {
								return err
							}
						}
						return nil
					}
					pterm.Warning.Println("Search index doesn't exist! Fetching the latest version.")
					err := indexer.DownloadIndex()
					if err != nil {
						return err
					}
					return nil
				},
				Action: func(c *cli.Context) error {
					results, err := indexer.LoadLocalIndex()
					if err != nil {
						panic(err)
					}
					_, err = search.FilterAlbumList(results)
					if err != nil {
						pterm.Info.Println("Goodbye")
					}
					return nil
				},
			},
			{
				Name:    "album",
				Aliases: []string{"a"},
				Usage:   "download an album given a slug",
				Action: func(c *cli.Context) error {
					tracks, err := scrape.DownloadAlbum(c.Args().First())
					if err != nil {
						panic(err)
					}
					fmt.Println(tracks)
					return nil
				},
			},
			{
				Name:    "index",
				Aliases: []string{"i"},
				Usage:   "builds an indexer of all khinsider content",
				Action: func(c *cli.Context) error {
					err := indexer.BuildIndex()
					if err != nil {
						panic(err)
					}
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
