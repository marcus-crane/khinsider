package khinsider

import (
	"errors"
	"github.com/marcus-crane/khinsider/v2/pkg/download"
	"github.com/marcus-crane/khinsider/v2/pkg/indexer"
	"github.com/marcus-crane/khinsider/v2/pkg/scrape"
	"github.com/marcus-crane/khinsider/v2/pkg/search"
	"github.com/marcus-crane/khinsider/v2/pkg/update"
	"github.com/pterm/pterm"
)

func BeforeSearch() error {
	indexExists := indexer.CheckIndexExists()
	if indexExists {
		pterm.Debug.Println("Checking for updates")
		localVersion := indexer.GetLocalIndexVersion()
		remoteVersion := update.GetRemoteIndexVersion()
		updateAvailable := update.IsRemoteVersionNewer(localVersion, remoteVersion)
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
}

func SearchAction() error {
	index, err := indexer.LoadLocalIndex()
	if err != nil {
		panic(err)
	}
	albumSlug, err := search.FilterAlbumList(index.Entries)
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
	album, err := scrape.RetrieveAlbum(albumSlug)
	if err != nil {
		panic(err)
	}
	download.GetAlbum(&album)
	return nil
}

func IndexAction() error {
	err := indexer.BuildIndex()
	if err != nil {
		panic(err)
	}
	return nil
}
