CREATE TABLE todo_records (
	id SERIAL PRIMARY KEY,
	title text NOT NULL,
	completed boolean NOT NULL,
	"order" integer NOT NULL
)
