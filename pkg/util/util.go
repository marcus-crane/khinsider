package util

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/pterm/pterm"
)

type Reader struct {
	io.Reader
	bar *pterm.ProgressbarPrinter
}

func NewBarProxyReader(r io.Reader, bar *pterm.ProgressbarPrinter) *Reader {
	return &Reader{r, bar}
}

func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	r.bar.Add(n)
	return
}

func (r *Reader) Close() (err error) {
	_, err = r.bar.Stop()
	if err != nil {
		return err
	}
	if closer, ok := r.Reader.(io.Closer); ok {
		return closer.Close()
	}
	return
}

func MakeRequest(link string, headers http.Header) (*http.Response, error) {
	headers.Add("User-Agent", "khinsider/3.0 <https://github.com/marcus-crane/khinsider>")
	remoteURL, err := url.Parse(link)
	if err != nil {
		panic(err)
	}
	client := http.Client{}
	request := http.Request{
		Method: "GET",
		URL:    remoteURL,
		Header: headers,
	}
	return client.Do(&request)
}

func RequestJSON(link string) (*http.Response, error) {
	headers := map[string][]string{
		"Accept-Encoding": {"application/json"},
		"Content-Type":    {"application/json"},
	}
	return MakeRequest(link, headers)
}

func RequestFile(link string) (*http.Response, error) {
	headers := map[string][]string{
		"Accept-Encoding": {"application/octet-stream"},
		"Content-Type":    {"application/octet-stream"},
	}
	return MakeRequest(link, headers)
}

func LoadJSON(file io.Reader, i interface{}) error {
	fileBytes, _ := io.ReadAll(file)
	err := json.Unmarshal(fileBytes, &i)
	if err != nil {
		pterm.Error.Println("Failed to load JSON")
		return errors.New("failed to load JSON")
	}
	return nil
}
