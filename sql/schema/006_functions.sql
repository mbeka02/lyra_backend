-- +goose Up
-- Function and trigger for appointment validation
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION check_appointment_availability()
RETURNS TRIGGER AS $BODY$
DECLARE
  slot_available boolean;
  has_overlap boolean;
BEGIN
  -- Check if the time slot exists in doctor's availability
  SELECT EXISTS (
    SELECT 1 FROM availability
    WHERE doctor_id = NEW.doctor_id
      AND (
        (is_recurring = true AND EXTRACT(DOW FROM NEW.start_time) = EXTRACT(DOW FROM start_time))
        OR
        (is_recurring = false AND specific_date = NEW.start_time::date)
      )
      AND start_time <= NEW.start_time::time
      AND end_time >= NEW.end_time::time
  ) INTO slot_available;

  IF NOT slot_available THEN
    RAISE EXCEPTION 'Time slot is not within doctor''s availability';
  END IF;

  -- Check for overlapping appointments
  SELECT EXISTS (
    SELECT 1 FROM appointments
    WHERE doctor_id = NEW.doctor_id
      AND current_status = 'scheduled'
      AND appointment_id != NEW.appointment_id
      AND (start_time, end_time) OVERLAPS (NEW.start_time, NEW.end_time)
  ) INTO has_overlap;

  IF has_overlap THEN
    RAISE EXCEPTION 'Time slot is already booked';
  END IF;

  RETURN NEW;
END;
$BODY$ LANGUAGE plpgsql;
-- +goose StatementEnd
-- Create trigger for appointment validation
CREATE TRIGGER validate_appointment
BEFORE INSERT OR UPDATE ON appointments
FOR EACH ROW
EXECUTE FUNCTION check_appointment_availability();

-- +goose Down
DROP TRIGGER IF EXISTS validate_appointment ON appointments;
DROP FUNCTION IF EXISTS check_appointment_availability();
