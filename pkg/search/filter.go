package search

import (
  //"github.com/manifoldco/promptui"

  "fmt"
  "sort"
  "strings"

  "github.com/manifoldco/promptui"
  "github.com/pterm/pterm"

  "github.com/marcus-crane/khinsider/pkg/types"
)

func FilterAlbumList(list types.SearchResults) (string, error) {
  keys := make([]string, 0, len(list))
  for k := range list {
    keys = append(keys, k)
  }
  sort.Strings(keys)
  prompt := promptui.Select{
    Label: "Search for an album",
    Items: keys,
    Size: 10,
    StartInSearchMode: true,
    Searcher: func(input string, index int) bool {
      album := keys[index]
      name := strings.Replace(strings.ToLower(album), " ", "", -1)
      input = strings.Replace(strings.ToLower(input), " ", "", -1)
      return strings.Contains(name, input)
    },
  }

  _, result, err := prompt.Run()

  pterm.Info.Printf("Selected %s\n", result)

  if err != nil {
    fmt.Printf("Prompt failed %v\n", err)
    return "", err
  }
  return result, nil
}
