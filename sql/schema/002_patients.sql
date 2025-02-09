-- +goose Up 
CREATE TABLE patients (
    patient_id SERIAL PRIMARY KEY,
    user_id INT UNIQUE NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    date_of_birth DATE,
    allergies TEXT
);
CREATE INDEX idx_patients_user_id ON patients(user_id);
-- +goose Down
DROP TABLE patients;
