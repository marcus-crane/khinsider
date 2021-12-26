package util

import (
	"io"

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
