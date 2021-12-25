package khinsider

import (
	"fmt"
	"log"
	"os"

	"github.com/marcus-crane/khinsider/pkg/scrape"
  "github.com/marcus-crane/khinsider/pkg/search"
  "github.com/urfave/cli/v2"
)

func Execute() {
	app := &cli.App{
		Name:    "khinsider",
		Usage:   "easily fetch videogame soundtracks from khinsider.com",
		Version: "2.0.0",
		Commands: []*cli.Command{
      {
        Name:    "search",
        Aliases: []string{"s"},
        Usage:   "search for an album to download",
        Action: func(c *cli.Context) error {
          results, err := scrape.GetResultsForLetter("A")
          if err != nil {
            panic(err)
          }
          _, err = search.FilterAlbumList(results)
          if err != nil {
            panic(err)
          }
          return nil
        },
      },
      {
        Name: "album",
        Aliases: []string{"a"},
        Usage: "download an album given a slug",
        Action: func(c *cli.Context) error {
          tracks, err := scrape.DownloadAlbum(c.Args().First())
          if err != nil {
            panic(err)
          }
          fmt.Println(tracks)
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
