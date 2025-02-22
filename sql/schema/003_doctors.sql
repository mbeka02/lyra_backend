-- +goose Up 
CREATE TABLE IF NOT EXISTS doctors(
    doctor_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    specialization VARCHAR(255) NOT NULL,
    license_number VARCHAR(50) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT (now()),
    updated_at timestamptz DEFAULT now()
  );
CREATE INDEX idx_doctors_user_id ON doctors(user_id);
CREATE INDEX idx_doctors_specialization ON doctors(specialization);
-- +goose Down
DROP TABLE doctors;
