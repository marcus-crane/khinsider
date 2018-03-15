package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/urfave/cli"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Album struct {
	Tracks []Track
}

type Track struct {
	Number int
	Title  string
	Link   string
}

func ScrapeAlbum(url string) *Album {
	resp, err := http.Get("https://downloads.khinsider.com/game-soundtracks/album/" + url)
	if err != nil {
		panic(err)
	}
	// No 404s, only 200 redirect to homepage
	if resp.StatusCode == 404 {
		return nil
	}
	root, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}

	album := Album{}

	matcher := func(n *html.Node) bool {
		// If the href contains ".mp3" and it doesn't have a 2nd attribute on td (align: right;)
		if n.DataAtom == atom.A && n.Parent.DataAtom == atom.Td && strings.Contains(n.Attr[0].Val, ".mp3") && len(n.Parent.Attr) == 1 {
			return scrape.Attr(n.Parent, "class") == "clickable-row"
		}
		return false
	}
	songs := scrape.FindAll(root, matcher)
	for i, song := range songs {
		track := Track{i + 1, scrape.Text(song), scrape.Attr(song, "href")}
		album.Tracks = append(album.Tracks, track)
	}
	return &album
}

func main() {
	app := cli.NewApp()
	app.Name = "khinsider"
	app.Usage = "Fetch albums from download.khinsider.com"
	app.Action = func(c *cli.Context) error {
		album := c.Args().Get(0)
		if album != "" {
			queryResults := ScrapeAlbum(album)
			if queryResults == nil {
				err := errors.New("sorry, that album doesn't seem to exist")
				panic(err)
			}
			for i := range queryResults.Tracks {
				track := queryResults.Tracks[i]
				fmt.Printf("%02d %s\n", track.Number, track.Title)
			}
			return nil
		}
		return errors.New("Please enter the name of an album eg katamari-damacy-soundtrack-katamari-fortissimo-damacy")
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Print(err)
	}
}
