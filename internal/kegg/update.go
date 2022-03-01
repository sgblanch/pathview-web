package kegg

import (
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/sgblanch/pathview-web/internal/config"
	"github.com/spf13/cobra"
)

func Update() {
	var (
		start = time.Now()
		c     = config.Get().Kegg

		err error
	)

	c.KeggDir, err = os.MkdirTemp(c.BaseDir, start.Format("20060102."))
	cobra.CheckErr(err)
	log.Printf("downloading to %q", c.KeggDir)

	d := Downloader{Directory: c.KeggDir}
	watch(&d, urls()...)
	// watch(&d, "http://rest.kegg.jp/list/organism")

	orgs, err := organisms(c.KeggDir)
	cobra.CheckErr(err)
	log.Printf("loading %d organisms", len(orgs))

	addl := flatten(map[string][]string{
		"http://rest.kegg.jp/list":         orgs,
		"http://rest.kegg.jp/link/pathway": orgs,
		"http://rest.kegg.jp/conv": flatten(map[string][]string{
			"ncbi-geneid":    orgs,
			"ncbi-proteinid": orgs,
			"uniprot":        orgs,
		}),
	})
	watch(&d, addl...)

	log.Print("creating database")
	createDB()

	end := time.Now()
	log.Printf("update finished in %v", end.Sub(start))
}

func watch(d *Downloader, urls ...string) {
	var (
		done = make(chan bool)
		wg   sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		d.Download(urls...)
		done <- true
		wg.Done()
	}()

	for {
		select {
		case <-done:
			fmt.Println()
			goto fin
		default:
			fmt.Printf("\r%d%% [%d:%d:%d] [completed:total:failed]", d.Progress(), d.Completed(), d.Items(), d.Failed())
			time.Sleep(1 * time.Second)
		}
	}
fin:
	wg.Wait()
}

func urls() []string {
	databases := []string{"compound", "disease", "drug", "genome", "glycan"}
	external := []string{"chebi", "pubchem"}
	conv := map[string][]string{
		"compound": external,
		"drug":     external,
		"glycan":   external,
	}
	urlMap := map[string][]string{
		"http://rest.kegg.jp/list":         append(databases, "organism", "pathway/ko", "pathway/map"),
		"http://rest.kegg.jp/link/pathway": databases,
		"http://rest.kegg.jp/conv":         flatten(conv),
	}

	return flatten(urlMap)
}

func createDB() {
	err := config.Get().DB.ExecAll(_create)
	cobra.CheckErr(err)
	defer func() {
		if err != nil {
			e := config.Get().DB.ExecAll(_fail)
			cobra.CheckErr(e)
		}
	}()

	log.Print("noop")
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
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		s = append(s, record[1])
	}

	if err != nil && err != io.EOF {
		return nil, err
	}

	return s, nil
}
