package download

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/marcus-crane/khinsider/v3/pkg/types"
	"github.com/marcus-crane/khinsider/v3/pkg/util"
	"github.com/pterm/pterm"
)

func GetAlbum(album *types.Album) {
	usrHome, _ := os.UserHomeDir()
	downloadFolder := fmt.Sprintf("%s/Downloads", usrHome) // TODO: Offer a configuration option
	normalisedSlug := strings.ReplaceAll(album.Slug, ".", "")
	directoryPath := fmt.Sprintf("%s/%s", downloadFolder, normalisedSlug)
	// TODO: This should be checked before download since it takes ages to get here
	_, err := os.Stat(directoryPath)
	if !errors.Is(err, fs.ErrNotExist) && err != nil {
		panic(err)
	}
	if os.IsExist(err) {
		err := os.Remove(directoryPath)
		if err != nil {
			pterm.Error.Printfln("A folder already exists at %s. Please remove it to continue.", directoryPath)
			os.Exit(1)
		}
	}
	err = os.Mkdir(directoryPath, 0755)
	if os.IsExist(err) && err != nil {
		pterm.Error.Printfln("A folder already exists at %s. Please remove it to continue.", directoryPath)
		os.Exit(1)
	}
	if os.IsNotExist(err) && err != nil {
		panic(err)
	}
	pterm.Success.Printfln("Successfully created %s", directoryPath)
	for _, track := range album.Tracks {
		trackFmt := track.Title
		// Some of the numbering can be quite bad on khinsider so we shouldn't
		// assume the track numbers are any good!
		padLength := len(fmt.Sprintf("%d", album.Total.Tracks))
		trackFmt = fmt.Sprintf("%0*d %s", padLength, track.TrackNumber, trackFmt)
		if track.DiscNumber != 0 {
			// TODO: Padding for 10+ discs
			trackFmt = fmt.Sprintf("%dx%s", track.DiscNumber, trackFmt)
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
	trackFile := fmt.Sprintf("%s/%s.mp3", saveLocation, normaliseFileName(fileName))
	pterm.Debug.Printfln("Downloading %s", track.SourceMP3)
	res, err := util.RequestFile(track.SourceMP3)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)
	if err != nil {
		pterm.Debug.Printfln("There was an error downloading %s", track.SourceMP3)
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
		pterm.Debug.Printfln(err.Error())
		return err
	}
	_, err = io.Copy(writer, res.Body)
	if err != nil {
		pterm.Debug.Printfln("There was an error writing %s", fileName)
		return err
	}
	return nil
}

func normaliseFileName(title string) string {
	// TODO: Code dump from v1. Should be reviewed again.
	if !utf8.ValidString(title) {
		pterm.Debug.Printfln("Invalid title: %s", title)
		validString := make([]rune, 0, len(title))
		for i, r := range title {
			if r == utf8.RuneError {
				_, size := utf8.DecodeRuneInString(title[i:])
				if size == 1 {
					continue
				}
			}
			validString = append(validString, r)
		}
		pterm.Debug.Printfln("Normalised title: %s", string(validString))
		return string(validString)
	}
	return title
}
