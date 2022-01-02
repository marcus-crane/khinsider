package util

import (
	"io"
	"net/http"
	"net/url"

	"github.com/pterm/pterm"
)

const (
	SiteBase = "https://downloads.khinsider.com"
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
		"User-Agent":      {"khinsider/2.0 <https://github.com/marcus-crane/khinsider>"},
	}
	return MakeRequest(link, headers)
}

func RequestFile(link string) (*http.Response, error) {
	headers := map[string][]string{
		"Accept-Encoding": {"application/octet-stream"},
		"Content-Type":    {"application/octet-stream"},
		"User-Agent":      {"khinsider/2.0 <https://github.com/marcus-crane/khinsider>"},
	}
	return MakeRequest(link, headers)
}
