package model

type Kegg interface {
	File() string
	FromRecord([]string) (*Kegg, error)
	Create(...Kegg) error
	Initialize() error
	Finalize() error
}
