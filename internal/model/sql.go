package model

var _sql = map[string]string{
	// "file_store": `insert into pathview.file (name, organism, file_format, file_content, owner) values (:name, :organism, :file_format, :file_content, :owner) returning id, created_at, updated_at;`,
	"file_store": `insert into pathview.file (name, organism) values (:name, :organism) returning id, created_at, updated_at;`,

	"organism_insert":  `insert into kegg_new.organism (id, code, name, common) values (:id, :code, :name, :common);`,
	"organism_hide":    `update kegg_new.organism set hidden = true where id not in (select distinct(org_id) from kegg.organism_pathway);`,
	"organism_default": `select id,code,name,common from kegg.organism where hidden = false order by seq limit 8;`,
	"organism_fts":     `select id,code,name,common from kegg.organism where tsv @@ to_tsquery('simple', :fts) and hidden = false limit 8;`,
	"organism_code":    `select code from kegg.organism where id = :org_id limit 1;`,

	"pathway_insert":  `insert into kegg_new.pathway (id, name) values (:id, :name);`,
	"pathway_hide":    `update kegg_new.pathway set hidden = true where 1100 <= id and id < 1300;`,
	"pathway_default": `select p.id, p.name from kegg.pathway as p join kegg.organism_pathway as op on p.id = op.path_id where op.org_id = :org_id order by p.seq limit 8;`,
	"pathway_fts":     `select p.id, p.name from kegg.pathway as p join kegg.organism_pathway as op on p.id = op.path_id where op.org_id = :org_id and p.tsv @@ to_tsquery('simple', :fts) limit 8;`,

	"organism_pathway_ko":     `insert into kegg_new.organism_pathway (org_id, path_id) select 0,id from kegg_new.pathway;`,
	"organism_pathway_insert": `insert into kegg_new.organism_pathway (org_id, path_id) values ((select id from kegg_new.organism where code = :code), :path_id);`,

	"user_create": `insert into pathview.user (name, email, email_verified, provider, provider_id) values (:name, :email, :email_verified, :provider, :provider_id) returning id, created_at, updated_at;`,
	"user_select": `select * from pathview.user where provider = :provider and provider_id = :provider_id;`,
	"user_update": `update pathview.user set (name, email, email_verified) = (:name, :email, :email_verified) where provider = :provider and provider_id = :providerid returning updated_at;`,
}
