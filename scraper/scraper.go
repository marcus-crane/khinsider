package scraper

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func ScrapeAlbum(url string) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	root, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}

	matcher := func(n *html.Node) bool {
		// If the href contains ".mp3" and it doesn't have a 2nd attribute on td (align: right;)
		if n.DataAtom == atom.A && n.Parent.DataAtom == atom.Td && strings.Contains(n.Attr[0].Val, ".mp3") && len(n.Parent.Attr) == 1 {
			return scrape.Attr(n.Parent, "class") == "clickable-row"
		}
		return false
	}
	songs := scrape.FindAll(root, matcher)
	for i, song := range songs {
		fmt.Printf("%2d %s (%s)\n", i+1, scrape.Text(song), scrape.Attr(song, "href"))
	}
}
