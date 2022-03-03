package util

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sync"
	"syscall"
)

type Downloader struct {
	Directory string
	urls      []string
	completed int
	items     int
	failed    int
	mu        sync.Mutex
	progMu    sync.Mutex
}

func (p *Downloader) Download(url ...string) {
	p.mu.Lock()
	p.urls = url
	p.mu.Unlock()
}

func (p *Downloader) Run() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.progMu.Lock()
	p.items = len(p.urls)
	p.progMu.Unlock()

	for _, u := range p.urls {
		err := p.get(u)
		if err != nil {
			log.Print(err)
			p.progMu.Lock()
			p.failed++
			p.progMu.Unlock()
		}

		p.progMu.Lock()
		p.completed++
		p.progMu.Unlock()
	}
}

func (p *Downloader) get(u string) error {
	var (
		response *http.Response
		outfile  string
		fh       *os.File
		gzfh     *gzip.Writer
		err      error
	)

	parsed, err := url.Parse(u)
	if err != nil {
		return err
	}

	for i := 0; i < 3; i++ {
		response, err = http.Get(u)
		if err != nil {
			continue
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusBadRequest || response.StatusCode == http.StatusNotFound {
			return fmt.Errorf("http status [%d] fetching %q", response.StatusCode, u)
		}
		if err != nil {
			log.Printf("downloading %q failed (retrying)", u)
			log.Printf("%v (ignored)", err)
		}
		if response.StatusCode == http.StatusOK {

			outfile = path.Join(p.Directory, fmt.Sprintf("%s.gz", parsed.Path))
			err = os.MkdirAll(filepath.Dir(outfile), 0700)
			if err != nil {
				goto failed
			}
			fh, err = os.Create(outfile)
			if err != nil {
				goto failed
			}
			defer fh.Close()

			gzfh = gzip.NewWriter(fh)
			defer gzfh.Close()

			_, err = io.Copy(gzfh, response.Body)
			if err != nil {
				goto failed
			}

			break
		}
		goto failed

	failed:
		if !errors.Is(err, syscall.ECONNRESET) {
			break
		}
		if fh != nil {
			fh.Close()
		}
		if gzfh != nil {
			gzfh.Close()
		}
		log.Printf("downloading %q failed (retrying)", u)
		log.Printf("%v (ignored)", err)
	}

	if err != nil {
		log.Printf("downloading %q failed, giving up", u)
		return err
	}

	return nil
}

func (p *Downloader) Progress() int {
	p.progMu.Lock()
	defer p.progMu.Unlock()

	if p.items == 0 {
		return 0
	}
	return 100 * p.completed / p.items
}

func (p *Downloader) Completed() int {
	p.progMu.Lock()
	defer p.progMu.Unlock()

	return p.completed
}

func (p *Downloader) Total() int {
	p.progMu.Lock()
	defer p.progMu.Unlock()

	return p.items
}

func (p *Downloader) Failed() int {
	p.progMu.Lock()
	defer p.progMu.Unlock()

	return p.failed
}
