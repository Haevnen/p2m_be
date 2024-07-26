// Use DBML to define your database structure
// Docs: https://dbml.dbdiagram.io/docs

enum user_type {
	fulltime
	freelancer
}

Table users {
id integer [primary key]
nick_name varchar(20)
email varchar(200)
password varchar(100)
type user_type
is_active bit
}

