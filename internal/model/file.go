package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/sgblanch/pathview-web/internal/config"
	"github.com/sgblanch/pathview-web/internal/util"
)

type FileFormat uint8

const (
	FileFormatCSV FileFormat = iota
	FileFormatTSV
)

func (p FileFormat) String() string {
	switch p {
	case FileFormatCSV:
		return "text/csv"
	case FileFormatTSV:
		return "text/tab-separated-values"
	default:
		return fmt.Sprintf("%d", int(p))
	}
}

func (p *FileFormat) UnmarshalJSON(data []byte) error {
	var s *string

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch strings.ToLower(*s) {
	case "text/csv":
		*p = FileFormatCSV
	case "csv":
		*p = FileFormatCSV
	case "text/tab-separated-values":
		*p = FileFormatTSV
	case "tsv":
		*p = FileFormatTSV
	default:
		return fmt.Errorf("invalid format: %q", *s)
	}

	return nil
}

type FileContent uint8

const (
	FileContentGene FileContent = iota
	FileContentCompound
)

func (p FileContent) String() string {
	switch p {
	case FileContentGene:
		return "gene"
	case FileContentCompound:
		return "compound"
	default:
		return fmt.Sprintf("%d", int(p))
	}
}

func (p *FileContent) UnmarshalJSON(data []byte) error {
	var s *string

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch strings.ToLower(*s) {
	case "gene":
		*p = FileContentGene
	case "compound":
		*p = FileContentCompound
	default:
		return fmt.Errorf("invalid content: %q", *s)
	}

	return nil
}

type FileUploadMetadata struct {
	ID   uuid.UUID `db:"id" json:"id"`
	Name string    `db:"name" json:"name"`
	// Format   FileFormat  `db:"file_format" json:"format" sql:"type:file_format"`
	// Content  FileContent `db:"file_content" json:"content" sql:"type:file_content"`
	Format    string    `db:"file_format" json:"format" sql:"type:file_format"`
	Content   string    `db:"file_content" json:"content" sql:"type:file_content"`
	Organism  string    `db:"organism" json:"organism"`
	Owner     User      `db:"owner" json:"owner"`
	Size      int64     `db:"size" json:"size"`
	Checksum  []byte    `db:"checksum" json:"checksum"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// func (p *FileUploadMetadata) Store() error {
// 	return config.Get().DB.NamedGet(p, _sql["file_store"], p)
// }

type FileUploadRequest struct {
	ID       uuid.UUID
	Name     string
	Format   FileFormat  `db:"file_format" json:"format" sql:"type:file_format"`
	Content  FileContent `db:"file_content" json:"content" sql:"type:file_content"`
	Organism string
	Owner    sql.NullInt64
	Session  string
	Token    string `db:"-"`
}

func (p *FileUploadRequest) StoreRequest() error {
	var (
		token   *string
		encoded []byte
		err     error
		c       = config.Get().RedisPool.Get()
	)
	defer c.Close()

	p.ID, err = uuid.NewV4()
	if err != nil {
		return err
	}

	token, err = util.RandomToken(32)
	if err != nil {
		return err
	}
	p.Token = *token

	encoded, err = json.Marshal(p)
	if err != nil {
		return err
	}

	_, err = c.Do("SET", p.Token, fmt.Sprintf("upload:request:%v", encoded), "EX", "300")

	return err
}
