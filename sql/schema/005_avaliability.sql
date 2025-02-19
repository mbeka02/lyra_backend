-- +goose Up 
CREATE TABLE IF NOT EXISTS availability(
availability_id bigserial PRIMARY KEY,
doctor_id bigint NOT NULL REFERENCES doctors(doctor_id),
start_time timestamptz NOT NULL,
end_time timestamptz NOT NULL
);

CREATE INDEX idx_availability_doctor_id ON availability(doctor_id);
-- +goose Down
DROP TABLE availability;
