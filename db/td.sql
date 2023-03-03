-- DATABASE NAME: utube

CREATE TABLE IF NOT EXISTS actor (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() at time zone 'utc'),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() at time zone 'utc'),

    CONSTRAINT actor_email_unique UNIQUE (email)
);

CREATE OR REPLACE FUNCTION update_actor_updated_at_column() RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- INSERT INTO actor (email, password, name) VALUES ('test_email', 'test_password', 'test_name') RETURNING *;