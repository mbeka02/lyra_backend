-- +goose Up 
CREATE TABLE IF NOT EXISTS availability(
availability_id bigserial PRIMARY KEY,
doctor_id bigint NOT NULL REFERENCES doctors(doctor_id) ON DELETE CASCADE,
start_time time NOT NULL,
end_time time NOT NULL,
is_recurring boolean NOT NULL DEFAULT true,
created_at timestamptz NOT NULL DEFAULT (now()),
--TODO: Add a trigger to update this before update
updated_at timestamptz DEFAULT (now()),
day_of_week integer NOT NULL CHECK (day_of_week between 0 AND 6),
interval_minutes integer NOT NULL DEFAULT 60,
CONSTRAINT valid_time_range CHECK (start_time<end_time),
CONSTRAINT unique_availability_slot UNIQUE(doctor_id,day_of_week,start_time,end_time)
);

CREATE INDEX idx_availability_doctor_id ON availability(doctor_id);
-- +goose Down
DROP TABLE availability;
