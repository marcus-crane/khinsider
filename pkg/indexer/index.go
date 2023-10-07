package indexer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/marcus-crane/khinsider/v3/pkg/update"
	"github.com/pterm/pterm"

	"github.com/marcus-crane/khinsider/v3/pkg/scrape"
	"github.com/marcus-crane/khinsider/v3/pkg/types"
	"github.com/marcus-crane/khinsider/v3/pkg/util"
)

const (
	LocalIndex  = "index.json"
	RemoteIndex = "https://raw.githubusercontent.com/marcus-crane/khinsider-index/main/index.json"
)

func getCachePath(filename string) string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return filepath.Join(usr.HomeDir, ".cache/khinsider/", filename)
}

func createPathIfNotExists(path string) {
	fullPath := filepath.Dir(path)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		err := os.MkdirAll(fullPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}

func CheckIndexExists() bool {
	indexPath := getCachePath(LocalIndex)
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func BuildIndex() error {
	results := make(types.SearchResults)
	letters := []string{"#"}
	for i := 'A'; i <= 'Z'; i++ {
		letter := fmt.Sprintf("%c", i)
		letters = append(letters, letter)
	}
	p, _ := pterm.DefaultProgressbar.WithTotal(len(letters)).WithTitle("Building indexer").WithRemoveWhenDone(true).Start()
	for _, letter := range letters {
		p.UpdateTitle("Downloading results for " + letter)
		page := 1
		letterResults, more, err := scrape.GetResultsForLetter(letter)
		for {
			if more {
				page += 1
				letterUrl := fmt.Sprintf("%s?page=%d", letter, page)
				p.UpdateTitle(fmt.Sprintf("~ Downloading Page %d of %s", page, letter))
				results, evenMore, err := scrape.GetResultsForLetter(letterUrl)
				if err != nil {
					panic(err)
				}
				for k, v := range results {
					letterResults[k] = v
				}
				if !evenMore {
					break
				}
				more = evenMore
			} else {
				break
			}
		}
		if err != nil {
			panic(err)
		}
		for k, v := range letterResults {
			results[k] = v
		}
		p.Increment()
	}
	index := types.SearchIndex{
		IndexVersion: IncrementIndexVersion(),
		Entries:      results,
	}
	err := SaveIndex(index)
	if err != nil {
		return err
	}
	return nil
}

func GetLocalIndexVersion() string {
	index, err := LoadLocalIndex()
	if err != nil {
		panic(err)
	}
	return update.ValidateIndexVersion(index.IndexVersion, "local")
}

func IncrementIndexVersion() string {
	remoteVersion := update.GetRemoteIndexVersion()
	semverBits := strings.Split(remoteVersion, ".")
	patch := semverBits[len(semverBits)-1]
	patchAsNumber, _ := strconv.Atoi(patch)
	patch = strconv.Itoa(patchAsNumber + 1)
	semverBits[len(semverBits)-1] = patch
	return strings.Join(semverBits, ".")
}

func DownloadIndex() error {
	pterm.Info.Printfln("Downloading latest search index. This lets you search all of khinsider locally.")
	res, err := util.RequestJSON(RemoteIndex)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)
	p, _ := pterm.DefaultProgressbar.WithTotal(int(res.ContentLength)).WithTitle("Downloading index").WithRemoveWhenDone(true).WithShowPercentage(true).Start()
	reader := util.NewBarProxyReader(res.Body, p)
	if res.StatusCode == http.StatusOK {
		pterm.Debug.Printfln("Retrieved index with status code of %d", res.StatusCode)
		indexPath := getCachePath(LocalIndex)
		createPathIfNotExists(indexPath)
		var index types.SearchIndex
		if err := util.LoadJSON(reader, &index); err != nil {
			panic(err)
		}
		err = SaveIndex(index)
		if err != nil {
			panic(err)
		}
		return nil
	}
	return fmt.Errorf("received a non-200 status code: %d", res.StatusCode)
}

func LoadLocalIndex() (types.SearchIndex, error) {
	file, err := os.Open(getCachePath(LocalIndex))
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	var index types.SearchIndex
	if err := util.LoadJSON(file, &index); err != nil {
		return index, nil
	}
	return index, err
}

func SaveIndex(index types.SearchIndex) error {
	if len(index.Entries) == 0 {
		pterm.Error.Println("It appears you have an empty index which shouldn't be possible. TODO: Build a bug report")
		return errors.New("index is empty which there is nothing to search through. that shouldn't be possible")
	}
	output, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		pterm.Error.Println("It appears you have a corrupted index! TODO: Build a bug report")
		return errors.New("index is malformed. can't serialise internal representation to json")
	}
	indexPath := getCachePath(LocalIndex)
	createPathIfNotExists(indexPath)
	err = os.WriteFile(indexPath, output, 0644)
	if err != nil {
		return errors.New("there was an issue writing the latest index to disc")
	}
	pterm.Debug.Printf("Successfully generated index at %s\n", indexPath)
	return nil
}
