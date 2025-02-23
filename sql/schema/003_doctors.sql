-- +goose Up
CREATE TABLE IF NOT EXISTS doctors(
    doctor_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    specialization VARCHAR(255) NOT NULL,
    years_of_experience INT NOT NULL,
    -- Maybe normalize this --
    county VARCHAR(30) NOT NULL,
    price_per_hour NUMERIC(10,2) NOT NULL,
    license_number VARCHAR(50) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT (now()),
    updated_at timestamptz DEFAULT now()
  );
CREATE INDEX idx_doctors_user_id ON doctors(user_id);
CREATE INDEX idx_doctors_specialization ON doctors(specialization);
CREATE INDEX idx_doctors_price ON doctors(price_per_hour);
CREATE INDEX idx_doctors_experience ON doctors(years_of_experience);
-- +goose Down
DROP TABLE doctors;
