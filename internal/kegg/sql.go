package kegg

var _sql = map[string]string{
	"file_insert":     `insert into kegg_new.file (file, sha256) values (:file, :sha256);`,
	"metadata_insert": `insert into kegg_new.metadata (key, value) values (:key, :value);`,
}
