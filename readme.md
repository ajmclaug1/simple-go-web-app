Postcres setup

create docker container runing postgres

docker run --name reading-db-container -e POSTGRES_PASSWORD=mysecretpassword -e POSTGRES_USER=postgres -d -p 5432:5432 postgres

Confirm its runninng

docker ps 

connect using psql

psql -h localhost -p 5432 -U postgres

password: mysecretpassword

SQL commands

CREATE DATABASE readinglist;

CREATE ROLE readinglist WITH LOGIN PASSWORD 'pa55w0rd';

\c readinglist;

CREATE TABLE IF NOT EXISTS books( 
	id bigserial PRIMARY KEY,
	created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
	title text NOT NULL, 
	published integer NOT NULL,
	pages integer NOT NULL,
	genres text[] NOT NULL,
	rating real NOT NULL,
	version integer NOT NULL DEFAULT 1
	);
	
GRANT SELECT, INSERT, UPDATE, DELETE ON books TO readinglist;

GRANT USAGE, SELECT ON SEQUENCE books_id_seq TO readinglist;
