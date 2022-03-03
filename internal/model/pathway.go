package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/sgblanch/pathview-web/internal/config"
	"github.com/sgblanch/pathview-web/internal/db"

	"github.com/gin-gonic/gin"
)

// Pathway foo
type Pathway struct {
	// Serial int       `db:"serial" json:"-"`
	ID   PathwayID `db:"id" json:"id"`
	Name string    `db:"name" json:"name"`
	// Hidden bool      `db:"hidden" json:"-"`
}

func (p Pathway) Default(organism int) ([]Pathway, error) {
	var (
		paths []Pathway
		err   error
	)
	err = config.Get().DB.NamedSelect(&paths, _sql["pathway_default"], gin.H{"org_id": organism})
	if err != nil {
		return nil, err
	}

	return paths, nil
}

func (p Pathway) Search(query string, organism int) ([]Pathway, error) {
	var (
		paths    []Pathway
		code     string
		ftsQuery string
	)

	err := config.Get().DB.NamedGet(&code, _sql["organism_code"], gin.H{"org_id": organism})
	if err != nil {
		return nil, err
	}

	ftsQuery = db.FtsQuery(query, code)

	if query == "" || ftsQuery == "" {
		return p.Default(organism)
	}

	err = config.Get().DB.NamedSelect(&paths, _sql["pathway_fts"], gin.H{"org_id": organism, "fts": ftsQuery})
	if err != nil {
		return nil, err
	}

	return paths, nil
}

func (p Pathway) File() string {
	return "list/pathway/ko.gz"
}

func (p Pathway) FromRecord(record []string) (*Kegg, error) {
	var pathway Kegg

	id, err := strconv.Atoi(strings.TrimPrefix(record[0], "path:ko"))
	if err != nil {
		return nil, err
	}

	pathway = Pathway{
		ID:   PathwayID(id),
		Name: record[1],
	}

	return &pathway, nil
}

func (p Pathway) Create(pathways ...Kegg) error {
	_, err := config.Get().DB.Chunk(240, _sql["pathway_insert"], pathways)

	return err
}

func (p Pathway) Initialize() error {
	return nil
}

func (p Pathway) Finalize() error {
	_, err := config.Get().DB.Exec(_sql["organism_pathway_hide"])
	return err
}

type PathwayID int

// MarshalJSON foo
func (p PathwayID) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%05d", p))
}

// UnmarshalJSON foo
func (p *PathwayID) UnmarshalJSON(data []byte) error {
	var s *string

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if s != nil && *s != "" {
		i, err := strconv.Atoi(*s)
		if err != nil {
			return err
		}

		*p = PathwayID(i)
	}

	return nil
}
