package kegg

import (
	"compress/gzip"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/sgblanch/pathview-web/internal/config"
	"github.com/sgblanch/pathview-web/internal/util"
)

func Update() {
	var (
		start = time.Now()
		c     = config.Get().Kegg

		// err error
	)

	c.KeggDir = path.Join(c.BaseDir, "20220303.2996394947")
	// c.KeggDir, err = os.MkdirTemp(c.BaseDir, start.Format("20060102."))
	// cobra.CheckErr(err)

	// log.Printf("downloading to %q", c.KeggDir)
	// download()

	log.Print("creating database")
	createDB()

	log.Print("loading database")
	loadDB()

	log.Print("indexing database")
	indexDB()

	end := time.Now()
	log.Printf("update finished in %v", end.Sub(start))
}

func watch(p util.Progresser) {
	var (
		done = make(chan bool)
		wg   sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		p.Run()
		done <- true
		wg.Done()
	}()

	for {
		select {
		case <-done:
			fmt.Println()
			goto fin
		default:
			fmt.Printf("\r%d%% [%d:%d:%d] [completed:total:failed]", p.Progress(), p.Completed(), p.Total(), p.Failed())
			time.Sleep(1 * time.Second)
		}
	}
fin:
	wg.Wait()
}

func organisms(file string) ([]string, error) {
	var (
		s      []string
		record []string
	)
	file = path.Join(file, "list/organism.gz")
	log.Printf("opening %q", file)
	fh, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	gzfh, err := gzip.NewReader(fh)
	if err != nil {
		return nil, err
	}
	defer gzfh.Close()

	r := csv.NewReader(gzfh)
	r.Comma = '\t'
	for {
		record, err = r.Read()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, err
		}
		s = append(s, record[1])
	}

	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}

	return s, nil
}
