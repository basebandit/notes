package mysql

var migrate = []string{
	`
		create table if not exists notes (
						id          bigint       not null auto_increment,
						name        varchar(60)  default null,
						author      varchar(60)  default null,
						created_at  datetime(6)  not null,
						deleted     boolean      not null default false,

						primary key (id),
						unique index(name),
						index(created_at)
		) default charset = utf8mb4;
	`,
}

var drop = []string{
	`drop table if exists notes`,
}
