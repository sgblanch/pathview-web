package db

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type DB struct {
	*sqlx.DB
	Host     string `mapstructure:"host"`
	Database string `mapstructure:"database"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

func (p *DB) Open() (err error) {
	p.DB, err = sqlx.Open("pgx", p.dsn())

	return
}

func (p *DB) dsn() string {
	var dsn strings.Builder

	if p.Host != "" {
		fmt.Fprintf(&dsn, "host=%v ", quote(p.Host))
	}
	if p.Database != "" {
		fmt.Fprintf(&dsn, "dbname=%v ", quote(p.Database))
	}
	if p.User != "" {
		fmt.Fprintf(&dsn, "user=%v ", quote(p.User))
	}
	if p.Password != "" {
		fmt.Fprintf(&dsn, "password=%v ", quote(p.Password))
	}
	dsn.WriteString("sslmode=disable")

	return dsn.String()
}

func (p *DB) ExecAll(queries []string) (err error) {
	var tx *sqlx.Tx

	tx, err = p.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			e := tx.Rollback()
			if e != nil {
				log.Printf("Rollback failed: %v", e)
			}
			return
		}
		err = tx.Commit()
	}()

	for _, v := range queries {
		_, err = tx.Exec(v)
		if err != nil {
			return fmt.Errorf("%q: %w", v, err)
		}
	}

	return
}

func (p *DB) NamedGet(dest interface{}, query string, arg interface{}) error {
	prepared, err := p.PrepareNamed(query)
	if err != nil {
		return err
	}

	return prepared.Get(dest, arg)
}

func (p *DB) NamedSelect(dest interface{}, query string, arg interface{}) error {
	prepared, err := p.PrepareNamed(query)
	if err != nil {
		return err
	}

	return prepared.Select(dest, arg)
}

func quote(str string) string {
	r := regexp.MustCompile(`('|\\)`)
	return r.ReplaceAllString(str, `\$1`)
}
