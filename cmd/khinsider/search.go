package khinsider

import (
	"errors"
	"strings"

	"github.com/marcus-crane/khinsider/v3/pkg/download"
	"github.com/marcus-crane/khinsider/v3/pkg/indexer"
	"github.com/marcus-crane/khinsider/v3/pkg/scrape"
	"github.com/marcus-crane/khinsider/v3/pkg/search"
	"github.com/pterm/pterm"
)

func BeforeSearch() error {
	pterm.Warning.Println("Fetching the latest search index.")
	err := indexer.DownloadIndex()
	if err != nil {
		return err
	}
	return nil
}

func SearchAction() error {
	index, err := indexer.LoadLocalIndex()
	if err != nil {
		panic(err)
	}
	albumSlug, err := search.FilterAlbumList(index)
	if err != nil {
		panic(err)
	}
	err = DownloadAction(albumSlug)
	if err != nil {
		pterm.Error.Println("Failed to download album")
	}
	return nil
}

func DownloadAction(albumSlug string) error {
	if albumSlug == "" {
		pterm.Error.Println("Please enter the slug for a valid album")
		return errors.New("no album slug provided")
	}
	// At present, the index captures entries as URL paths so eg;
	// /game-soundtrack/album/<slug> whereas the user downloads
	// and album by providing just the slug. We could update the
	// index to just save slugs but this would break compatibility
	// with earlier versions so instead we just strip the index
	// entries down to their slug. Both searching and direct slug
	// download pass through this function so they need to be consistent
	if strings.Contains(albumSlug, "/") {
		slugBits := strings.Split(albumSlug, "/")
		albumSlug = slugBits[len(slugBits)-1]
	}
	album, err := scrape.RetrieveAlbum(albumSlug)
	if err != nil {
		panic(err)
	}
	download.GetAlbum(&album)
	return nil
}
