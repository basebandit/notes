package postgresql

var migrate = []string{
	`     create table if not exists notes (
								id          bigserial    not null primary key,
								name        text         default null,
								author      text         default null,
								content     text         not null,
								created_at  timestamptz  not null,
								deleted     boolean      not null default false
				);
				create unique index on notes(lower(name));
				create unique index on notes(created_at);`,
}

var drop = []string{
	`drop table if exists notes`,
}
