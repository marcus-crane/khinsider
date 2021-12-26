package download

import (
	"fmt"

	"github.com/marcus-crane/khinsider/pkg/types"
)

func GetAlbum(album *types.Album) {
	fmt.Printf("Downloading %s", album.Name)
}
