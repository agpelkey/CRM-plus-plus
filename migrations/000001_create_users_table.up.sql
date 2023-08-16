CREATE TABLE IF NOT EXISTS users (
	id bigserial PRIMARY KEY,
	first_name text NOT NULL,
	last_name text NOT NULL,
	phone_number text NOT NULL,
	email text NOT NULL,
	follow_up bool NOT NULL DEFAULT FALSE,
	check_in_date date NOT NULL DEFAULT NOW()
);

