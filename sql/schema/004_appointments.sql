-- +goose Up 
CREATE TYPE appointment_status AS ENUM ('pending_payment','scheduled', 'completed', 'canceled');
CREATE TABLE IF NOT EXISTS appointments(
  appointment_id bigserial PRIMARY KEY,
  patient_id bigint NOT NULL REFERENCES patients(patient_id) ON DELETE CASCADE,
  doctor_id bigint NOT NULL REFERENCES doctors(doctor_id) ON DELETE CASCADE,
  current_status appointment_status NOT NULL DEFAULT ('scheduled'),
  reason TEXT NOT NULL,
notes TEXT,
  start_time timestamptz NOT NULL,
  end_time timestamptz NOT NULL,
  created_at timestamptz NOT NULL DEFAULT (now()),
  updated_at timestamptz DEFAULT (now()),
  CONSTRAINT valid_appointment_time CHECK (start_time<end_time)
);
CREATE INDEX idx_appointments_patient_id ON appointments(patient_id);
CREATE INDEX idx_appointments_doctor_id ON appointments(doctor_id);
-- +goose Down
DROP TABLE appointments;
DROP TYPE appointment_status;

