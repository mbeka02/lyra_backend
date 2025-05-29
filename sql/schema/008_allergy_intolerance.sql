-- +goose Up

CREATE TABLE IF NOT EXISTS allergy_intolerances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    patient_id BIGINT NOT NULL REFERENCES patients(patient_id) ON DELETE CASCADE,
    clinical_status_code VARCHAR(50) NOT NULL,       -- e.g., 'active', 'inactive', 'resolved'
    clinical_status_display TEXT,
    code_system VARCHAR(255),                        -- System for the allergy code
    code_code VARCHAR(100) NOT NULL,                 -- The allergy code itself
    code_display TEXT NOT NULL,                      -- Human-readable display of the allergy
    criticality VARCHAR(50),                         -- e.g., 'low', 'high', 'unable-to-assess' (nullable)
    reaction_manifestation_text TEXT,                -- Simplified: Store manifestation as text. (e.g., "Rash, Itching")
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_allergy_intolerances_patient_id ON allergy_intolerances(patient_id);

-- +goose Down
DROP TABLE allergy_intolerances;
