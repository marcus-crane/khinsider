package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"strings"
	"unicode/utf8"

	"github.com/urfave/cli/v2"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"gopkg.in/cheggaaa/pb.v2"
)

type Album struct {
	Tracks []Track
}

type Track struct {
	Number int
	Title  string
	Link   string
}

func scrapeAlbum(url string) *Album {
	resp, err := http.Get("https://downloads.khinsider.com/game-soundtracks/album/" + url)
	if err != nil {
		panic(err)
	}
	// No 404s, only 200 redirect to homepage
	if resp.StatusCode == http.StatusNotFound {
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

func pullAudioStream(url string) string {
	resp, err := http.Get("https://downloads.khinsider.com" + url)
	if err != nil {
		panic(err)
	}
	root, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}

	matcher := func(n *html.Node) bool {
		if n.DataAtom == atom.Audio {
			return true
		}
		return false
	}
	stream, _ := scrape.Find(root, matcher)
	return scrape.Attr(stream, "src")
}

func downloadFile(filepath string, url string) (err error) {
	writer, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer writer.Close()

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Uh oh, failed to fetch: %s", resp.Status)
	}
	defer resp.Body.Close()

	bar := pb.ProgressBarTemplate(`{{ counters . }} {{ bar . "[" "=" (rnd "" ) "-" "]"}} {{speed . }}`).Start64(resp.ContentLength)
	bar.Start()

	reader := bar.NewProxyReader(resp.Body)

	io.Copy(writer, reader)
	bar.Finish()
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "khinsider"
	app.Usage = "Fetch albums from download.khinsider.com"
	app.Version = "1.11.0" // damn, versioning for this sucks huh
	app.Action = func(c *cli.Context) error {
		album := c.Args().Get(0)
		if album != "" {
			queryResults := scrapeAlbum(album)
			if queryResults == nil {
				err := errors.New("sorry, that album doesn't seem to exist")
				panic(err)
			}
			usr, err := user.Current()
			if err != nil {
				panic(err)
			}
			os.Mkdir(usr.HomeDir+"/Downloads/"+album, 0755)
			fmt.Printf("Created %s\n\n", usr.HomeDir+"/Downloads/"+album)
			for i := range queryResults.Tracks {
				track := queryResults.Tracks[i]
				track.Link = pullAudioStream(track.Link)
				fmt.Printf("Downloading %02d %s\n", track.Number, track.Title)
				// this should hold for now until i do a proper rewrite
				if !utf8.ValidString(track.Title) {
					validString := make([]rune, 0, len(track.Title))
					for i, r := range track.Title {
						if r == utf8.RuneError {
							_, size := utf8.DecodeRuneInString(track.Title[i:])
							if size == 1 {
								continue
							}
						}
						validString = append(validString, r)
					}
					track.Title = string(validString)
				}
				filePath := fmt.Sprintf(usr.HomeDir+"/Downloads/%s/%02d %s.mp3", album, track.Number, track.Title)
				downloadFile(filePath, track.Link)
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
