// https://dbml.dbdiagram.io/docs/

enum contract_type {
	fulltime
	freelancer
}

Table users {
	id integer [increment, primary key]
    user_id char(36) [not null, unique]
	nick_name varchar(20) [not null, unique]
	email varchar(200) [not null, unique]
	password_hashed varchar(255) [not null]
	contract_type contract_type
	is_active bit [not null, default: `1`]
    is_admin bit [not null, default: `0`]
	created_at timestamp [not null, default: `now()`]
}

Table sessions {
	id bigint [increment, primary key]
    session_id char(36) [not null, unique]
	user_id char(36) [not null, ref: > users.user_id]
	refresh_token varchar(255) [not null, unique]
	created_at timestamp [not null, default: `now()`]
	expired_at timestamp [not null]
}