package kegg

var _create = []string{
	`drop schema if exists kegg_new cascade;`,

	`create schema kegg_new;`,

	`create table kegg_new.metadata (
		key varchar(32) primary key,
		value varchar(128)
	);`,

	`create table kegg_new.organism (
		seq serial,
		id integer primary key,
		code varchar(5) unique not null,
		uniprot varchar(8),
		ncbi_tax_id integer, 
		name varchar(126) not null,
		common varchar(126),
		hidden boolean not null default false
	);`,

	// `create table kegg_new.genome (
	// 	id integer primary key,
	// 	code varchar(5) unique,
	// 	uniprot varchar(5),
	// 	ncbi_tax_id integer
	// );`,

	`create table kegg_new.pathway (
		seq serial,
		id integer primary key,
		name varchar(128) unique not null,
		hidden boolean not null default false
	);`,

	`create table kegg_new.organism_pathway (
		org_id integer not null,
		path_id integer not null,
		primary key(org_id, path_id),
		foreign key(org_id) references kegg_new.organism(id),
		foreign key(path_id) references kegg_new.pathway(id)
	);`,

	`create table kegg_new.gene (
		id character varying(64) not null,
		org_id integer not null,
		name character varying(512) not null,
		primary key(id, org_id),
		foreign key(org_id) references kegg_new.organism(id)
	);`,

	`create table kegg_new.gene_alias (
		gene_id character varying(64) not null,
		org_id integer not null,
		alias character varying(128) not null,
		primary key(gene_id, alias),
		foreign key(gene_id, org_id) references kegg_new.gene(id, org_id)
	);`,

	`create table kegg_new.file (
		file varchar(128) primary key,
		sha256 char(64) not null
	);`,
}

var _index = []string{
	`create unique index on kegg_new.organism(code);`,
	`create index on kegg_new.organism_pathway(org_id);`,
	// `update kegg_new.organism as o set (uniprot, ncbi_tax_id) = (select uniprot, ncbi_tax_id from kegg_new.genome as g where g.code = o.code);`,
	`alter table kegg_new.organism add column tsv tsvector generated always as (to_tsvector('simple', id::text || ' ' || code || ' ' || coalesce(uniprot, '') || ' ' || coalesce(ncbi_tax_id::text, '') || ' ' || name || ' ' || coalesce(common, ''))) stored;`,
	`create index on kegg_new.organism using gin (tsv);`,
	`alter table kegg_new.pathway add column tsv tsvector generated always as (to_tsvector('simple', id::text || ' ' || name)) stored;`,
	`create index on kegg_new.pathway using gin (tsv);`,
	`drop table if exists kegg_new.genome;`,
}

var _fail = []string{
	`drop schema if exists kegg_new cascade;`,
}

// var _drop = []string{
// 	`drop schema if exists kegg_old cascade;`,
// 	`drop schema if exists kegg cascade;`,
// 	`drop schema if exists kegg_new cascade;`,
// }

var _pivot = []string{
	`drop schema if exists kegg_old cascade;`,
	`do $$begin if exists(select 1 from information_schema.schemata where schema_name = 'kegg') then execute 'alter schema kegg rename to kegg_old'; end if; end $$;`,
	`alter schema kegg_new rename to kegg;`,
}
