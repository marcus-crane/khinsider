package indexer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"github.com/pterm/pterm"

	"github.com/marcus-crane/khinsider/pkg/scrape"
	"github.com/marcus-crane/khinsider/pkg/types"
	"github.com/marcus-crane/khinsider/pkg/util"
)

const (
	ETagLocation    = "cache_etag.txt"
	IndexLocation   = "index.json"
	RemoteIndexGist = "https://gist.githubusercontent.com/marcus-crane/7b27c3b2772115ce7dbb7be792a7ad75/raw/index.json"
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
	indexPath := getConfigPath(IndexLocation)
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
	err := SaveIndex(results)
	if err != nil {
		return err
	}
	return nil
}

func CheckIndexUpdateAvailable() bool {
	pterm.Debug.Printfln("Fetching index headers from %s", RemoteIndexGist)
	res, err := requestIndex("HEAD")
	if err != nil {
		panic(err)
	}
	pterm.Debug.Println("Successfully fetched headers")
	etag := res.Header.Get("Etag")
	if etag == "" {
		pterm.Debug.Println("No ETag header was found on the remote index!")
		panic(err)
	}
	lastEtag := loadLastEtag()
	pterm.Debug.Printfln("Loaded current ETag from disk. Current version is %s", lastEtag)
	pterm.Debug.Printfln("Fetched ETag header from request. Remote version is %s", etag)
	if etag != lastEtag {
		pterm.Debug.Printfln("The two versions don't match! Assuming the remote contains a new update.")
		return true
	}
	pterm.Debug.Println("The two versions match! We already have the latest index on disc.")
	return false
}

func requestIndex(method string) (*http.Response, error) {
	url, _ := url.Parse(RemoteIndexGist)
	client := http.Client{}
	request := http.Request{
		Method: "GET",
		URL:    url,
		Header: map[string][]string{
			"Accept-Encoding": {"application/json"},
			"Content-Type":    {"application/json"},
		},
	}
	return client.Do(&request)
}

func DownloadIndex() error {
	pterm.Info.Printfln("Downloading latest search index. This lets you search all of khinsider locally.")
	res, err := requestIndex("GET")
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
	etag := res.Header.Get("Etag")
	pterm.Debug.Printfln("ETag for downloaded index is %s", etag)
	if res.StatusCode == http.StatusOK {
		pterm.Debug.Printfln("Retrieved index with status code of %d", res.StatusCode)
		indexPath := getConfigPath(IndexLocation)
		createPathIfNotExists(indexPath)
		results, err := LoadIndex(reader)
		if err != nil {
			panic(err)
		}
		err = SaveIndex(results)
		if err != nil {
			panic(err)
		}
		err = SaveETag(etag)
		if err != nil {
			panic(err)
		}
		return nil
	}
	return fmt.Errorf("received a non-200 status code: %d", res.StatusCode)
}

func loadLastEtag() string {
	eTagStore := getConfigPath(ETagLocation)
	if _, err := os.Stat(eTagStore); os.IsNotExist(err) {
		return ""
	}
	contents, err := ioutil.ReadFile(eTagStore)
	if err != nil {
		return ""
	}
	return string(contents)
}

func LoadLocalIndex() (types.SearchResults, error) {
	file, err := os.Open(getConfigPath(IndexLocation))
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	return LoadIndex(file)
}

func LoadIndex(file io.Reader) (types.SearchResults, error) {
	results := make(types.SearchResults)
	fileBytes, _ := ioutil.ReadAll(file)
	err := json.Unmarshal(fileBytes, &results)
	if err != nil {
		pterm.Error.Println("It appears you have a corrupted index! TODO: Build a bug report")
		return nil, errors.New("index is malformed. can't serialise internal representation to json")
	}
	return results, nil
}

func SaveETag(tag string) error {
	eTagFile := getConfigPath(ETagLocation)
	createPathIfNotExists(eTagFile)
	err := os.WriteFile(eTagFile, []byte(tag), 0644)
	if err != nil {
		return errors.New("there was an issue writing etag to disc")
	}
	pterm.Debug.Println("Saved ETag to disk")
	return nil
}

func SaveIndex(results types.SearchResults) error {
	if len(results) == 0 {
		pterm.Error.Println("It appears you have an empty index which shouldn't be possible. TODO: Build a bug report")
		return errors.New("index is empty. search functionality won't work with no results")
	}
	output, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		pterm.Error.Println("It appears you have a corrupted index! TODO: Build a bug report")
		return errors.New("index is malformed. can't serialise internal representation to json")
	}
	indexPath := getConfigPath(IndexLocation)
	createPathIfNotExists(indexPath)
	err = os.WriteFile(indexPath, output, 0644)
	if err != nil {
		return errors.New("there was an issue writing the latest index to disc")
	}
	pterm.Debug.Printf("Successfully generated index at %s\n", indexPath)
	return nil
}
