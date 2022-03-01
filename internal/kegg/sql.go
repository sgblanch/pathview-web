package kegg

var _sql = map[string]string{
	"file_insert":             `insert into kegg_new.file (file, sha256) values (:file, :sha256);`,
	"metadata_insert":         `insert into kegg_new.metadata (key, value) values (:key, :value);`,
	"organism_insert":         `insert into kegg_new.organism (id, code, name, common) values (:id, :code, :name, :common);`,
	"organism_select":         `select id, code from kegg_new.organism;`,
	"organism_hide":           `update kegg_new.organism set hidden = true where id = ?;`,
	"organism_id":             `select id from kegg_new.organism where code = ? limit 1;`,
	"organism_pathway_insert": `insert into kegg_new.organism_pathway (org_id, path_id) values (:org_id, :path_id);`,
	"organism_pathway_select": `select path_id as id from kegg_new.organism_pathway where org_id = ?;`,
	"pathway_insert":          `insert into kegg_new.pathway (id, name) values (:id, :name);`,
	"pathway_hide":            `update kegg_new.pathway set hidden = true where where 1100 <= id and id < 1300;`,
	"pathway_ko":              `select 0 as org_id, id as path_id from kegg_new.pathway;`,
}
