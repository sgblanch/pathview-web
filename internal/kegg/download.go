package kegg

import (
	"fmt"
	"log"

	"github.com/sgblanch/pathview-web/internal/config"
	"github.com/sgblanch/pathview-web/internal/util"
	"github.com/spf13/cobra"
)

func download() {
	var c = config.Get().Kegg

	d := util.Downloader{Directory: c.KeggDir}
	d.Download(urls()...)
	watch(&d)

	orgs, err := organisms(c.KeggDir)
	cobra.CheckErr(err)
	log.Printf("loading %d organisms", len(orgs))

	addl := flatten(map[string][]string{
		"http://rest.kegg.jp/list": orgs,
		// Obviated by /link/pathway/genomes
		// "http://rest.kegg.jp/link/pathway": orgs,
		// These take *forever*
		"http://rest.kegg.jp/conv": flatten(map[string][]string{
			"ncbi-geneid":    orgs,
			"ncbi-proteinid": orgs,
			"uniprot":        orgs,
		}),
	})
	d.Download(addl...)
	watch(&d)
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

func flatten(s map[string][]string) []string {
	var urls []string

	for url, fragments := range s {
		for _, fragment := range fragments {
			urls = append(urls, fmt.Sprintf("%s/%s", url, fragment))
		}
	}

	return urls
}
