-- +goose Up
CREATE TABLE IF NOT EXISTS observations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),     -- Or BIGSERIAL
    patient_id BIGINT NOT NULL REFERENCES patients(patient_id) ON DELETE CASCADE,
    -- specialist_id BIGINT REFERENCES doctors(doctor_id) ON DELETE SET NULL, -- Optional: if doctor enters note
    status VARCHAR(50) NOT NULL,                       -- e.g., 'final', 'amended'
    code_text TEXT NOT NULL,                           -- Description of what was observed (e.g., "Consultation Note")
    effective_date_time TIMESTAMPTZ NOT NULL,
    value_string TEXT NOT NULL,                        -- The actual note content
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_observations_patient_id ON observations(patient_id);

-- +goose Down
DROP TABLE observations;
