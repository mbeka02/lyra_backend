-- +goose Up 
CREATE TABLE patients (
    patient_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    date_of_birth DATE NOT NULL,
    allergies TEXT NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT (now()),
    updated_at timestamptz DEFAULT now()
);
CREATE INDEX idx_patients_user_id ON patients(user_id);
-- +goose Down
DROP TABLE patients;
