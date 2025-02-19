-- +goose Up 
CREATE TYPE status AS ENUM ('scheduled', 'completed', 'canceled');
CREATE TABLE IF NOT EXISTS appointments(
  appointment_id bigserial PRIMARY KEY,
  patient_id bigint NOT NULL REFERENCES patients(patient_id),
  doctor_id bigint NOT NULL REFERENCES doctors(doctor_id),
  current_status status NOT NULL DEFAULT ('scheduled'),
  appointment_date timestamptz NOT NULL
);
CREATE INDEX idx_appointments_patient_id ON appointments(patient_id);
CREATE INDEX idx_appointments_doctor_id ON appointments(doctor_id);
CREATE INDEX idx_appointments_date_status ON appointments(appointment_date, current_status);

-- +goose Down
DROP TABLE appointments;
DROP TYPE status;

