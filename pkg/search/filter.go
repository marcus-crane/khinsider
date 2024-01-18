package search

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/marcus-crane/khinsider/v3/pkg/scrape"
	"github.com/marcus-crane/khinsider/v3/pkg/types"
)

type item struct {
	Title string
	Meta  types.AlbumHints
}

func secondsToMinutes(inSeconds int32) string {
	minutes := inSeconds / 60
	seconds := inSeconds % 60
	str := fmt.Sprintf("%d:%d", minutes, seconds)
	return str
}

func FilterAlbumList(list types.SearchResults) ([]string, error) {
	filterList := []item{}
	for k, v := range list {
		filterList = append(filterList, item{Title: k, Meta: v})
	}
	sort.SliceStable(filterList, func(i, j int) bool {
		return strings.ToLower(filterList[i].Title) < strings.ToLower(filterList[j].Title)
	})
	idx, err := fuzzyfinder.FindMulti(
		filterList,
		func(i int) string {
			return filterList[i].Title
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			album, err := scrape.RetrieveAlbum(filterList[i].Meta.Slug)
			if err != nil {
				return "Failed to retrieve album metadata"
			}
			trackList := ""
			for _, track := range album.Tracks {
				trackList += fmt.Sprintf("%d. %s (%s)\n", track.TrackNumber, track.Title, secondsToMinutes(track.Runtime))
			}
			return fmt.Sprintf(`%s (%d)
			MP3 Available: %t
			FLAC Available: %t
			Genre: %s
			Track Count: %d
			Disc Count: %d

			Tracks:
			%s
			`,
				filterList[i].Title,
				filterList[i].Meta.Year,
				filterList[i].Meta.MP3Exists,
				filterList[i].Meta.FlacExists,
				filterList[i].Meta.Genre,
				filterList[i].Meta.TrackCount,
				filterList[i].Meta.DiscCount,
				trackList)
		}))
	if err != nil {
		log.Fatal(err)
	}
	selectedSlugs := []string{}
	for _, i := range idx {
		selectedSlugs = append(selectedSlugs, filterList[i].Meta.Slug)
	}
	return selectedSlugs, nil
}
