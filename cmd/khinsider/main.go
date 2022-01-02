package khinsider

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"

	"github.com/marcus-crane/khinsider/v2/pkg/download"
	"github.com/marcus-crane/khinsider/v2/pkg/indexer"
	"github.com/marcus-crane/khinsider/v2/pkg/scrape"
	"github.com/marcus-crane/khinsider/v2/pkg/search"
)

type BuildInfo struct {
	Version string
	Commit  string
	Date    string
	BuiltBy string
}

func (b BuildInfo) Print() {
	fmt.Println("Version:\t", b.Version)
	fmt.Println("Commit Hash:\t", b.Commit)
	fmt.Println("Build Date:\t", b.Date)
	fmt.Println("Build Source:\t", b.BuiltBy)
}

func Execute(buildInfo BuildInfo) {
	cli.VersionPrinter = func(c *cli.Context) {
		buildInfo.Print()
	}
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
		Usage: "easily fetch videogame soundtracks from downloads.khinsider.com",
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
						updateAvailable := indexer.IsRemoteVersionNewer()
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
					index, err := indexer.LoadLocalIndex()
					if err != nil {
						panic(err)
					}
					_, err = search.FilterAlbumList(index.Entries)
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
					albumSlug := c.Args().First()
					if albumSlug == "" {
						pterm.Error.Println("Please enter the slug for a valid album")
						return errors.New("no album slug provided")
					}
					album, err := scrape.DownloadAlbum(albumSlug)
					if err != nil {
						panic(err)
					}
					download.GetAlbum(&album)
					return nil
				},
			},
			{
				Name:    "index",
				Aliases: []string{"i"},
				Usage:   "generates a local index of all khinsider content",
				Hidden:  true,
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
