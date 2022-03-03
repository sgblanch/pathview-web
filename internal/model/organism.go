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

type Organism struct {
	ID     OrganismID    `json:"id"`
	Code   string        `json:"code"`
	Name   string        `json:"name"`
	Common db.NullString `json:"common,omitempty"`
	Tax    string        `json:"-" db:"-"`
	Hidden bool          `json:"-"`
}

func (p Organism) Default() ([]Organism, error) {
	var (
		orgs []Organism
		err  error
	)

	err = config.Get().DB.Select(&orgs, _sql["organism_default"])
	if err != nil {
		return nil, err
	}

	return orgs, nil
}

func (p Organism) Search(query string) ([]Organism, error) {
	var (
		orgs     []Organism
		ftsQuery = db.FtsQuery(query, "T")
	)

	if query == "" || ftsQuery == "" {
		return p.Default()
	}

	err := config.Get().DB.NamedSelect(&orgs, _sql["organism_fts"], gin.H{"fts": ftsQuery})
	if err != nil {
		return nil, err
	}

	return orgs, nil
}

func (p Organism) File() string {
	return "list/organism.gz"
}

func (p Organism) FromRecord(record []string) (*Kegg, error) {
	var (
		id     int
		common db.NullString
		err    error
	)

	id, err = strconv.Atoi(strings.TrimPrefix(record[0], "T"))
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(record[3], "Eukaryotes;") {
		start := strings.Index(record[2], " (")
		end := strings.Index(record[2], ")")
		if start >= 0 && end >= 0 {
			common = db.NullString{
				String: record[2][start+2 : end],
				Valid:  true,
			}
			if end < len(record[2]) {
				record[2] = record[2][:start] + record[2][end+1:]
			} else {
				record[2] = record[2][:start]
			}
		} else {
			common = db.NullString{Valid: false}
		}
	}

	var organism Kegg
	organism = Organism{
		ID:     OrganismID(id),
		Code:   record[1],
		Name:   record[2],
		Common: common,
		Tax:    record[3],
	}

	return &organism, nil
}

func (p Organism) Create(organisms ...Kegg) error {
	_, err := config.Get().DB.Chunk(240, _sql["organism_insert"], organisms)

	return err
}

func (p Organism) Initialize() error {
	return p.Create(Organism{
		ID:   0,
		Code: "ko",
		Name: "Kegg Orthology",
		Common: db.NullString{
			String: "Model Pathways",
			Valid:  true,
		},
	})
}

func (p Organism) Finalize() error {
	_, err := config.Get().DB.Exec(_sql["organism_hide"])
	return err
}

func (p *Organism) MarshalJSON() ([]byte, error) {
	type Alias Organism
	var common *string

	if p.Common.Valid {
		common = &p.Common.String
	}

	return json.Marshal(&struct {
		*Alias
		Common *string `json:"common,omitempty"`
	}{
		Alias:  (*Alias)(p),
		Common: common,
	})
}

type OrganismID int

// MarshalJSON foo
func (p OrganismID) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("T%05d", p))
}

// UnmarshalJSON foo
func (p *OrganismID) UnmarshalJSON(data []byte) error {
	var s *string

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if s != nil && *s != "" {
		i, err := strconv.Atoi(strings.TrimPrefix(*s, "T"))
		if err != nil {
			return err
		}

		*p = OrganismID(i)
	}

	return nil
}
