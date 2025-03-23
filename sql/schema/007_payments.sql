-- +goose Up
CREATE TYPE payment_status AS ENUM ('pending', 'completed','failed');
CREATE TABLE IF NOT EXISTS payments(
payment_id BIGSERIAL PRIMARY KEY,
reference VARCHAR UNIQUE NOT NULL,
current_status payment_status NOT NULL DEFAULT('pending'),
amount NUMERIC(10,2) NOT NULL,
metadata JSONB,
payment_method VARCHAR(30),
currency VARCHAR(4) NOT NULL DEFAULT('KSH'),
appointment_id BIGINT UNIQUE NOT NULL references appointments(appointment_id),
patient_id BIGINT NOT NULL references patients(patient_id),
doctor_id BIGINT NOT NULL references doctors(doctor_id),
created_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
updated_at TIMESTAMPTZ DEFAULT (now()),
completed_at TIMESTAMPTZ
);
CREATE UNIQUE INDEX idx_payments_appointment_id ON payments(appointment_id);
CREATE INDEX idx_payments_patient_id ON payments(patient_id);
CREATE INDEX idx_doctor_id ON payments(doctor_id);

-- +goose Down
DROP TABLE payments;
DROP TYPE payment_status;

