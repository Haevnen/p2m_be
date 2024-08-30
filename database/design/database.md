// https://dbml.dbdiagram.io/docs/

enum contract_type {
	fulltime
	freelancer
}

enum ticket_status {
    backlog
    in_progress
    ready_to_qc
    qc_verifying
    qc_done
    done
}

enum ticket_priority {
    normal
    high
}

enum ticket_created_by {
    auto
    manual
}

enum ticket_action {
    status_changed
}


Table users {
	id integer [increment, primary key]
    user_id char(36) [not null, unique]
	nick_name varchar(20) [not null, unique]
	email varchar(200) [not null, unique]
	password_hashed varchar(255) [not null]
	contract_type contract_type
	is_active tinyint(1) [not null, default: `1`]
    is_admin tinyint(1) [not null, default: `0`]
    is_unassigned tinyint(1) [not null, default: `0`]
	created_at timestamp [not null, default: `now()`]
}

Table sessions {
	id bigint [increment, primary key]
    session_id char(36) [not null, unique]
	user_id char(36) [not null, ref: > users.user_id]
	refresh_token varchar(512) [not null, unique]
	created_at timestamp [not null, default: `now()`]
	expired_at timestamp [not null]
}

Table clients {
    id integer [increment, primary key]
    client_id varchar(20)
    editing_style text
    requirements text
    other_info text
    is_active tinyint(1) [not null, default: `1`]
    created_at timestamp [not null, default: `now()`]
}

Table tickets {
    id bigint [increment, primary key]
    title varchar(255) [not null]
    status ticket_status [not null]
    qc_id char(36) [not null, ref: > users.user_id]
    editor_id char(36) [not null, ref: > users.user_id]
    priority ticket_priority [not null]
    client_id integer [not null, ref: > clients.id]
    description text
    is_active tinyint(1) [not null, default: `1`]
    created_by ticket_created_by [not null]
    created_at timestamp [not null, default: `now()`]
    updated_at timestamp [not null, default: `now()`]
}

Table links {
    id bigint [increment, primary key]
    ticket_id bigint [not null, ref: > tickets.id]
    link varchar(2083) [not null]
    created_at timestamp [not null, default: `now()`]
}

Table comments {
    id bigint [increment, primary key]
    ticket_id bigint [not null, ref: > tickets.id]
    user_id char(36) [not null, ref: > users.user_id]
    comment text [not null]
    created_at timestamp [not null, default: `now()`]
}

Table histories {
    id bigint [increment, primary key]
    ticket_id bigint [not null, ref: > tickets.id]
    action text 
    performed_by char(36) [not null, ref: > users.user_id]
    created_at timestamp [not null, default: `now()`]
}