package search

import (
	"fmt"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/marcus-crane/khinsider/v3/pkg/types"
)

func FilterAlbumList(list types.SearchResults) (types.AlbumHints, error) {
	keys := make([]string, 0, len(list))
	for k, v := range list {
		kFmt := fmt.Sprintf("%s ~ %d", k, v.Year)
		if v.DiscCount > 0 {
			kFmt = fmt.Sprintf("%s | D%d", kFmt, v.DiscCount)
		}
		kFmt = fmt.Sprintf("%s | T%d | %s", kFmt, v.TrackCount, v.Genre)
		if v.MP3Exists {
			kFmt = fmt.Sprintf("%s | [MP3]", kFmt)
		}
		if v.FlacExists {
			kFmt = fmt.Sprintf("%s [FLAC]", kFmt)
		}
		keys = append(keys, kFmt)
	}
	sort.Strings(keys)

	prompt := &survey.Select{
		Message: "Choose an album:",
		Options: keys,
	}

	var result string
	err := survey.AskOne(prompt, &result, survey.WithPageSize(25))

	if err != nil {
		return types.AlbumHints{}, err
	}
	title := strings.Split(result, " ~")[0]
	return list[title], nil
}
