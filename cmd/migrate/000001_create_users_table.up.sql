CREATE TABLE IF NOT EXISTS users(
	id bigserial NOT NULL,
	first_name text NOT NULL,
	last_name text NOT NULL,
	phone_number text NOT NULL,
	email text NOT NULL,
	created_at timestamptz DEFAULT NOW(),

	PRIMARY KEY(id)
);
