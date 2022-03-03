package kegg

import (
	"compress/gzip"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"path"

	"github.com/sgblanch/pathview-web/internal/config"
	"github.com/sgblanch/pathview-web/internal/model"
	"github.com/spf13/cobra"
)

func createDB() {
	err := config.Get().DB.ExecAll(_create)
	cobra.CheckErr(err)
	defer func() {
		if err != nil {
			e := config.Get().DB.ExecAll(_fail)
			cobra.CheckErr(e)
		}
	}()
}

func loadDB() {
	var m = []model.Kegg{
		model.Organism{},
		model.Pathway{},
		model.OrgPathway{},
	}

	for _, v := range m {
		err := v.Initialize()
		cobra.CheckErr(err)
		loadTable(v)
	}

	for _, v := range m {
		err := v.Finalize()
		cobra.CheckErr(err)
	}
}

func loadTable(m model.Kegg) {
	var (
		record []string
		row    *model.Kegg
		rows   []model.Kegg
		file   = path.Join(config.Get().Kegg.KeggDir, m.File())
	)

	fh, err := os.Open(file)
	cobra.CheckErr(err)
	defer fh.Close()

	gzfh, err := gzip.NewReader(fh)
	cobra.CheckErr(err)
	defer gzfh.Close()

	r := csv.NewReader(gzfh)
	r.Comma = '\t'
	for {
		record, err = r.Read()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			cobra.CheckErr(err)
		}

		row, err = m.FromRecord(record)
		cobra.CheckErr(err)

		rows = append(rows, *row)
	}

	err = m.Create(rows...)
	cobra.CheckErr(err)
}

func indexDB() {
	err := config.Get().DB.ExecAll(_index)
	cobra.CheckErr(err)

	err = config.Get().DB.ExecAll(_pivot)
	cobra.CheckErr(err)
}
