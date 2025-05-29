-- +goose Up
CREATE TABLE IF NOT EXISTS medication_statements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),   
    patient_id BIGINT NOT NULL REFERENCES patients(patient_id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL,                     -- e.g., 'active', 'completed', 'stopped'
    medication_code_system VARCHAR(255),
    medication_code_code VARCHAR(100) NOT NULL,
    medication_code_display TEXT NOT NULL,
    dosage_text TEXT NOT NULL,                                -- Simplified dosage instructions as text
    effective_date_time TIMESTAMPTZ,                 -- Can be NULL if not specified or ongoing
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_medication_statements_patient_id ON medication_statements(patient_id);

-- +goose Down
DROP TABLE medication_statements;
