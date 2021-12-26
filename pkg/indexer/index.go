package indexer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pterm/pterm"

	"github.com/marcus-crane/khinsider/pkg/scrape"
	"github.com/marcus-crane/khinsider/pkg/types"
)

const (
	IndexLocation = "$HOME/.khinsider/index.json"
)

func createIndexPathIfNotExists() {
	path := filepath.Dir(IndexLocation)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}

func BuildIndex() error {
	results := make(types.SearchResults)
	letters := []string{"#"}
	for i := 'A'; i <= 'Z'; i++ {
		letter := fmt.Sprintf("%c", i)
		letters = append(letters, letter)
	}
	p, _ := pterm.DefaultProgressbar.WithTotal(len(letters)).WithTitle("Building indexer").Start()
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
	output, err := json.Marshal(results)
	if err != nil {
		panic(err)
	}
	createIndexPathIfNotExists()
	err = os.WriteFile(IndexLocation, output, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("successfully generated index")
	return nil
}

func LoadIndex() (types.SearchResults, error) {
	results := make(types.SearchResults)
	file, err := os.Open(IndexLocation)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	fileBytes, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(fileBytes, &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}
