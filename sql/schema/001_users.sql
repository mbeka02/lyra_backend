-- +goose Up
CREATE TYPE role AS ENUM ('patient', 'specialist');
CREATE TABLE IF NOT EXISTS users (
user_id  bigserial PRIMARY KEY,
date_of_birth DATE NOT NULL,
full_name varchar(256) NOT NULL,
password varchar NOT NULL,
email varchar UNIQUE NOT NULL,
telephone_number varchar(16) NOT NULL,
profile_image_url text NOT NULL DEFAULT '' ,
created_at timestamptz NOT NULL DEFAULT (now()),
user_role role NOT NULL,
verified_at timestamptz,
is_onboarded boolean NOT NULL DEFAULT false,
password_changed_at timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

CREATE  UNIQUE INDEX idx_users_email ON users(email);
CREATE INDEX idx_user_role ON users(user_role);
-- +goose Down
DROP TABLE users;
DROP TYPE role;

