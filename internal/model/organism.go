package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/sgblanch/pathview-web/internal/config"
	"github.com/sgblanch/pathview-web/internal/db"
	dbx "github.com/sgblanch/pathview-web/internal/db"

	"github.com/gin-gonic/gin"
)

type Organism struct {
	ID     OrganismID    `json:"id"`
	Code   string        `json:"code"`
	Name   string        `json:"name"`
	Common db.NullString `json:"common,omitempty"`
	// Tax    string        `json:"-" db:"-"`
	// Hidden bool          `json:"-"`
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
		ftsQuery = dbx.FtsQuery(query, "T")
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
