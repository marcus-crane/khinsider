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

	"github.com/pterm/pterm"

	"github.com/marcus-crane/khinsider/v3/pkg/types"
	"github.com/marcus-crane/khinsider/v3/pkg/util"
)

const (
	LocalIndex  = "index-v3.json"
	RemoteIndex = "https://khindex.utf9k.net/index.json"
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
	if res.StatusCode == http.StatusOK {
		pterm.Debug.Printfln("Retrieved index with status code of %d", res.StatusCode)
		indexPath := getCachePath(LocalIndex)
		createPathIfNotExists(indexPath)
		var index types.SearchResults
		if err := util.LoadJSON(res.Body, &index); err != nil {
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

func LoadLocalIndex() (types.SearchResults, error) {
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
	var index types.SearchResults
	if err := util.LoadJSON(file, &index); err != nil {
		return index, nil
	}
	return index, err
}

func SaveIndex(index types.SearchResults) error {
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
