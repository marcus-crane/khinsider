package scrape

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"

	"github.com/marcus-crane/khinsider/pkg/types"
)

const (
	AlbumBase  = "https://downloads.khinsider.com/game-soundtracks/album/"
	LetterBase = "https://downloads.khinsider.com/game-soundtracks/browse/"
)

func DownloadPage(url string) (*http.Response, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, err
	}
	return res, err
}

func GetResultsForLetter(letter string) (types.SearchResults, error) {
	url := fmt.Sprintf("%s%s", LetterBase, letter)
	res, err := DownloadPage(url)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	results := make(types.SearchResults)
	doc.Find("#EchoTopic p[align='left'] a").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		results[title] = "#"
		url, exists := s.Attr("href")
		if exists {
			results[title] = url
		} else {
			results[title] = "#"
		}
	})
	return results, nil
}

func DownloadAlbum(slug string) (types.Album, error) {
	var album types.Album
	url := fmt.Sprintf("%s%s", AlbumBase, slug)
	res, err := DownloadPage(url)
	defer res.Body.Close()
	if err != nil {
		return album, err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return album, err
	}
	metadata := doc.Find("#EchoTopic p[align='left'] b")
	if metadata.Length() == 5 {
		album.FlacAvailable = true
	}
	metadata.Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			album.Name = s.Text()
		}
		if i == 1 {
			album.FileCount, err = strconv.Atoi(s.Text())
			if err != nil {
				album.FileCount = 0
			}
		}
		if i == 2 {
			album.MP3FileSize = s.Text()
		}
		if i == 3 && album.FlacAvailable {
			album.FlacFileSize = s.Text()
		}
	})
	doc.Find("#songlist .clickable-row:not([align])").Each(func(i int, s *goquery.Selection) {
		//fmt.Println(s.Text())
	})
	return album, nil
}
