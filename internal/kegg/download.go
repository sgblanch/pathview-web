package kegg

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sync"
)

type Progresser interface {
	Progress() int
	Completed() int
	Items() int
}

type Failer interface {
	Failed() int
}

type Downloader struct {
	Directory string
	completed int
	items     int
	failed    int
	mu        sync.Mutex
}

func (p *Downloader) Download(url ...string) {
	p.mu.Lock()
	p.items = len(url)
	p.mu.Unlock()

	for _, u := range url {
		err := p.get(u)
		if err != nil {
			log.Print(err)
			p.mu.Lock()
			p.failed++
			p.mu.Unlock()
		}

		p.mu.Lock()
		p.completed++
		p.mu.Unlock()
	}
}

func (p *Downloader) get(u string) error {
	var (
		response *http.Response
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

		if response.StatusCode == http.StatusOK {
			break
		}
		if response.StatusCode == http.StatusBadRequest || response.StatusCode == http.StatusNotFound {
			return fmt.Errorf("http status [%d] fetching %q", response.StatusCode, u)
		}
	}
	if err != nil {
		return err
	}

	outfile := path.Join(p.Directory, fmt.Sprintf("%s.gz", parsed.Path))
	err = os.MkdirAll(filepath.Dir(outfile), 0700)
	if err != nil {
		return err
	}
	fh, err := os.Create(outfile)
	if err != nil {
		return err
	}
	defer fh.Close()

	gzfh := gzip.NewWriter(fh)
	defer gzfh.Close()

	_, err = io.Copy(gzfh, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func (p *Downloader) Progress() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.items == 0 {
		return 0
	}
	return 100 * p.completed / p.items
}

func (p *Downloader) Completed() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.completed
}

func (p *Downloader) Items() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.items
}

func (p *Downloader) Failed() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.failed
}

func flatten(s map[string][]string) []string {
	var urls []string

	for url, fragments := range s {
		for _, fragment := range fragments {
			urls = append(urls, fmt.Sprintf("%s/%s", url, fragment))
		}
	}

	return urls
}
