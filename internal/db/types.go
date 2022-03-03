package db

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type NullString sql.NullString

func (p *NullString) Scan(value interface{}) error {
	if value == nil {
		p.String, p.Valid = "", false
		return nil
	}
	p.Valid = true

	switch s := value.(type) {
	case string:
		p.String = s
		return nil
	default:
		return fmt.Errorf("scan: expected a string")
	}
}

func (p NullString) Value() (driver.Value, error) {
	if !p.Valid {
		return nil, nil
	}
	return p.String, nil
}

// MarshalJSON foo
func (p *NullString) MarshalJSON() ([]byte, error) {
	if p.Valid {
		return json.Marshal(p.String)
	}

	return json.Marshal(nil)
}

// UnmarshalJSON foo
func (p *NullString) UnmarshalJSON(data []byte) error {
	var s *string

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if s != nil {
		p.Valid = true
		p.String = *s
	}

	return nil
}
