package scrape

import (
	"fmt"
	"io"
	"net/http"

	"github.com/marcus-crane/khinsider/v3/pkg/util"

	"github.com/pterm/pterm"

	"github.com/marcus-crane/khinsider/v3/pkg/types"
)

const (
	IndexAlbumBase = "https://khindex.utf9k.net/albums"
)

func DownloadPage(url string) (*http.Response, error) {
	res, err := util.MakeRequest(url, http.Header{})
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received a non-200 status code: %d", res.StatusCode)
	}
	return res, err
}

func RetrieveAlbum(slug string) (types.Album, error) {
	var album types.Album
	albumUrl := fmt.Sprintf("%s/%s.json", IndexAlbumBase, slug)

	res, err := DownloadPage(albumUrl)
	if err != nil {
		return album, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)

	err = util.LoadJSON(res.Body, &album)
	if err != nil {
		return album, err
	}
	pterm.Success.Printfln("Retrieved metadata for %s (%d tracks)", album.Title, album.Total.Tracks)
	return album, nil
}
