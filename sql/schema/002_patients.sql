-- +goose Up 
CREATE TABLE patients (
    patient_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
   -- Demographics
    address TEXT NOT NULL DEFAULT '',
    emergency_contact_name VARCHAR(255) NOT NULL DEFAULT '',
    emergency_contact_phone VARCHAR(20) NOT NULL DEFAULT '',
 -- Medical Information
   -- blood_type VARCHAR(3),
   -- height_cm INT,
   -- weight_kg DECIMAL(5,2),
    allergies TEXT NOT NULL DEFAULT '',
    current_medication TEXT NOT NULL DEFAULT '',
    past_medical_history TEXT NOT NULL DEFAULT '',
    family_medical_history TEXT NOT NULL DEFAULT '',
    -- Insurance Information
    insurance_provider VARCHAR(255) NOT NULL DEFAULT '',
    insurance_policy_number VARCHAR(255) NOT NULL DEFAULT '',

    created_at timestamptz NOT NULL DEFAULT (now()),
    updated_at timestamptz DEFAULT now()
);
CREATE INDEX idx_patients_user_id ON patients(user_id);
-- +goose Down
DROP TABLE patients;
