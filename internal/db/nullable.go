package db

import (
	"database/sql"
	"encoding/json"
)

type NullString struct {
	sql.NullString
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
