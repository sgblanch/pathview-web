package model

import (
	"strconv"
	"strings"

	"github.com/sgblanch/pathview-web/internal/config"
)

type OrgPathway struct {
	Code   string    `db:"code"`
	PathID PathwayID `db:"path_id"`
}

func (p OrgPathway) File() string {
	return "link/pathway/genome.gz"
}

func (p OrgPathway) FromRecord(record []string) (*Kegg, error) {
	code := strings.TrimPrefix(record[0], "gn:")
	id, err := strconv.Atoi(strings.TrimPrefix(record[1], "path:"+code))
	if err != nil {
		return nil, err
	}

	var orgpathway Kegg
	orgpathway = OrgPathway{
		Code:   code,
		PathID: PathwayID(id),
	}

	return &orgpathway, nil
}

func (p OrgPathway) Create(organisms ...Kegg) error {
	// This query is "slow" (60s) as it does ~7500 sub-selects.  No good
	// way to optimize since record creation is abstracted which limits
	// caching of organism codes -> ids.  Probably not necessary as this
	// shouldn't be run often
	_, err := config.Get().DB.Chunk(240, _sql["organism_pathway_insert"], organisms)

	return err
}

func (p OrgPathway) Initialize() error {
	_, err := config.Get().DB.Exec(_sql["organism_pathway_ko"])
	return err
}

func (p OrgPathway) Finalize() error {
	return nil
}
