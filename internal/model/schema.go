package model

var _create = []string{
	`create extension if not exists citext;`,

	`create type status as enum ('pending', 'succeeded', 'failed');`,

	`create type file_format as enum ('text/csv', 'text/tab-separated-values');`,

	`create type file_content as enum ('gene', 'compound');`,

	`create schema pathview;`,

	`create table pathview.user (
		id bigserial primary key,
		name varchar(255),
		email citext unique,
		email_verified boolean default false,
		provider varchar(64) not null,
		provider_id varchar(255) not null,
		created_at timestamp with time zone not null default now(),
		updated_at timestamp with time zone not null default now()
	);`,

	`create table pathview.file (
		id uuid primary key default gen_random_uuid(),
		name varchar(255),
		file_format file_format,
		file_content file_content,
		organism varchar(5) references kegg.organism(code),
		owner bigint references pathview.user(id),
		size bigint not null default 0,
		checksum bytea not null,
		created_at timestamp with time zone not null default now(),
		updated_at timestamp with time zone not null default now()
	);`,

	`create table pathview.analysis (
		id uuid primary key default gen_random_uuid(),
		owner bigint references pathview.user(id),
		organism varchar(5) references kegg.organism(code),
		gene uuid references pathview.file(id),
		compound uuid references pathview.file(id),
		config jsonb not null,
		result jsonb,
		status status default 'pending',
		created_at timestamp with time zone not null default now(),
		updated_at timestamp with time zone not null default now()
	);`,
}

var _index = []string{
	`create unique index on pathview.user(provider, provider_id);`,

	`create or replace function pathview.update_timestamp()
	returns trigger as $$
	begin
		new.updated_at = now();
		return new;
	end;
	$$ language plpgsql;`,

	`create trigger update_timestamp
	before update on pathview.user
	for each row
	execute procedure update_timestamp();`,

	`create trigger update_timestamp
	before update on pathview.file
	for each row
	execute procedure update_timestamp();`,

	`create trigger update_timestamp
	before update on pathview.analysis
	for each row
	execute procedure update_timestamp();`,
}

var _seed = []string{
	`insert into pathview.user (name, email, provider, provider_id) values 
		('Public', 'public', 'internal', 'public@internal'),
		('Session', 'session', 'internal', 'session@internal');`,
}
