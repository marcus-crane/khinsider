package search

import (
	"fmt"
	"sort"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pterm/pterm"

	"github.com/marcus-crane/khinsider/pkg/types"
)

func FilterAlbumList(list types.SearchResults) (string, error) {
	keys := make([]string, 0, len(list))
	for k := range list {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	prompt := &survey.Select{
		Message: "Choose an album:",
		Options: keys,
	}

	var result string
	err := survey.AskOne(prompt, &result, survey.WithPageSize(15))

	pterm.Info.Printf("Selected %s\n", result)

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}
	return result, nil
}
