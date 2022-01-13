package download

import (
	"fmt"
	"github.com/marcus-crane/khinsider/v2/pkg/types"
	"github.com/marcus-crane/khinsider/v2/pkg/util"
	"github.com/pterm/pterm"
	"io"
	"os"
	"strings"
)

func GetAlbum(album *types.Album) {
	usrHome, _ := os.UserHomeDir()
	downloadFolder := fmt.Sprintf("%s/Downloads", usrHome) // TODO: Offer a configuration option
	normalisedSlug := strings.ReplaceAll(album.Slug, ".", "")
	directoryPath := fmt.Sprintf("%s/%s", downloadFolder, normalisedSlug)
	err := os.Remove(directoryPath)
	if err != nil {
		panic(err)
	}
	err = os.Mkdir(directoryPath, 0755)
	if err != nil {
		panic(err)
	}
	pterm.Success.Printfln("Successfully created %s", directoryPath)
	lastCDNumber := ""
	trackNumEndLastCD := 0
	for i := 0; i < album.FileCount; i++ {
		track := album.Tracks[i]
		trackFmt := track.Name
		// When we hit a new CD, we'll reset the numbering and start again from
		// zero. We need to make sure we don't set the wrong padding though.
		if track.CDNumber != lastCDNumber {
			lastCDNumber = track.CDNumber
			trackNumEndLastCD = i
		}
		// Some of the numbering can be quite bad on khinsider so we shouldn't
		// assume the track numbers are any good!
		padLength := len(fmt.Sprintf("%d", album.FileCount))
		if track.CDNumber == lastCDNumber {
			padLength = 2 // Assume most CDs aren't any bigger than 100 tracks
		}
		formatFormat := "%0" + fmt.Sprintf("%d", padLength) + "d %s"
		trackCount := i + 1
		if trackNumEndLastCD != 0 {
			trackCount = trackCount - trackNumEndLastCD
		}
		trackFmt = fmt.Sprintf(formatFormat, trackCount, trackFmt)
		if track.CDNumber != "" {
			trackFmt = fmt.Sprintf("%sx%s", track.CDNumber, trackFmt)
		}
		err := SaveAudioFile(track, trackFmt, directoryPath)
		if err != nil {
			pterm.Error.Printfln(trackFmt)
		} else {
			pterm.Success.Printfln(trackFmt)
		}
	}
	fmt.Println()

}

func SaveAudioFile(track types.Track, fileName string, saveLocation string) error {
	trackFile := fmt.Sprintf("%s/%s.mp3", saveLocation, fileName)
	pterm.Debug.Printfln("Downloading %s", track.URL)
	res, err := util.RequestFile(track.URL)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)
	if err != nil {
		pterm.Debug.Printfln("There was an error downloading %s", track.URL)
		return err
	}
	writer, err := os.Create(trackFile)
	defer func(writer *os.File) {
		err := writer.Close()
		if err != nil {
			panic(err)
		}
	}(writer)
	if err != nil {
		pterm.Debug.Printfln("There was an error creating %s", trackFile)
		return err
	}
	_, err = io.Copy(writer, res.Body)
	if err != nil {
		pterm.Debug.Printfln("There was an error writing %s", fileName)
		return err
	}
	return nil
}
