create table accounts
(
	id int auto_increment
		primary key,
	email varchar(320) not null,
	phone int null,
	password varchar(255) not null,
	constraint account_email_uindex
		unique (email)
);

