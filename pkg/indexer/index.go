package indexer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pterm/pterm"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/mod/semver"

	"github.com/marcus-crane/khinsider/v2/pkg/scrape"
	"github.com/marcus-crane/khinsider/v2/pkg/types"
	"github.com/marcus-crane/khinsider/v2/pkg/util"
)

const (
	LocalIndex       = "index.json"
	IndexReleaseFeed = "https://api.github.com/repos/marcus-crane/khinsider-index/releases/latest"
	RemoteIndex      = "https://raw.githubusercontent.com/marcus-crane/khinsider-index/main/index.json"
)

func getConfigPath(filename string) string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return filepath.Join(usr.HomeDir, ".khinsider", filename)
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
	indexPath := getConfigPath(LocalIndex)
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
		letterResults, err := scrape.GetResultsForLetter(letter)
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
	return ValidateIndexVersion(index.IndexVersion, "local")
}

func GetRemoteIndexVersion() string {
	res, err := util.RequestJSON(IndexReleaseFeed)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)
	var remoteIndex types.RemoteIndexMetadata
	if err := LoadJSON(res.Body, &remoteIndex); err != nil {
		panic(err)
	}
	return ValidateIndexVersion(remoteIndex.Version, "remote")
}

func ValidateIndexVersion(version string, indexLocation string) string {
	if !strings.HasPrefix(version, "v") {
		pterm.Error.Printfln("%s index version %s doesn't start with a v.", indexLocation, version)
		panic(errors.New(fmt.Sprintf("%s index version is invalid", indexLocation)))
	}
	versionBits := strings.Split(version, ".")
	if len(versionBits) != 3 {
		pterm.Error.Printf("expected %s version %s to have 3 parts. %s only has %d", indexLocation, version, len(versionBits))
		panic(errors.New(fmt.Sprintf("%s index version is invalid", indexLocation)))
	}
	return version
}

func IncrementIndexVersion() string {
	remoteVersion := GetRemoteIndexVersion()
	semverBits := strings.Split(remoteVersion, ".")
	patch := semverBits[len(semverBits)-1]
	patchAsNumber, _ := strconv.Atoi(patch)
	patch = strconv.Itoa(patchAsNumber + 1)
	semverBits[len(semverBits)-1] = patch
	return strings.Join(semverBits, ".")
}

func IsRemoteVersionNewer() bool {
	localVersion := GetLocalIndexVersion()
	remoteVersion := GetRemoteIndexVersion()
	result := semver.Compare(localVersion, remoteVersion)
	if result == -1 {
		return true
	}
	return false
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
		indexPath := getConfigPath(LocalIndex)
		createPathIfNotExists(indexPath)
		var index types.SearchIndex
		if err := LoadJSON(reader, &index); err != nil {
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

func LoadJSON(file io.Reader, i interface{}) error {
	fileBytes, _ := ioutil.ReadAll(file)
	err := json.Unmarshal(fileBytes, &i)
	if err != nil {
		pterm.Error.Println("Failed to load JSON")
		return errors.New("failed to load JSON")
	}
	return nil
}

func LoadLocalIndex() (types.SearchIndex, error) {
	file, err := os.Open(getConfigPath(LocalIndex))
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
	if err := LoadJSON(file, &index); err != nil {
		return index, nil
	}
	return index, err
}

func SaveIndex(index types.SearchIndex) error {
	if len(index.Entries) == 0 {
		pterm.Error.Println("It appears you have an empty index which shouldn't be possible. TODO: Build a bug report")
		return errors.New("index is empty. search functionality won't work with no results")
	}
	output, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		pterm.Error.Println("It appears you have a corrupted index! TODO: Build a bug report")
		return errors.New("index is malformed. can't serialise internal representation to json")
	}
	indexPath := getConfigPath(LocalIndex)
	createPathIfNotExists(indexPath)
	err = os.WriteFile(indexPath, output, 0644)
	if err != nil {
		return errors.New("there was an issue writing the latest index to disc")
	}
	pterm.Debug.Printf("Successfully generated index at %s\n", indexPath)
	return nil
}
