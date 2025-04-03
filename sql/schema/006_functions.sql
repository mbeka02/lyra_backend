
-- Function and trigger for appointment validation
-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION check_appointment_availability()
RETURNS TRIGGER AS $BODY$
DECLARE
  slot_available boolean;
  has_overlap boolean;
  appt_start_time time;
  appt_end_time time;
  appt_date date;
  appt_dow integer;
BEGIN
  -- Extract the time and date components from appointment timestamptz
  appt_start_time := (NEW.start_time)::time;
  appt_end_time := (NEW.end_time)::time;
  appt_dow := EXTRACT(DOW FROM NEW.start_time);
  
  -- Check if the time slot exists in doctor's availability
  SELECT EXISTS (
    SELECT 1 FROM availability
    WHERE doctor_id = NEW.doctor_id
      AND (is_recurring = true AND day_of_week=appt_dow)
      AND (start_time, end_time) OVERLAPS (appt_start_time, appt_end_time)
  ) INTO slot_available;
  
  IF NOT slot_available THEN
    RAISE EXCEPTION 'Time slot is not within doctor''s availability';
  END IF;
  
  -- Check for overlapping appointments
  SELECT EXISTS (
    SELECT 1 FROM appointments
    WHERE doctor_id = NEW.doctor_id
      AND current_status = 'scheduled'
      AND appointment_id != COALESCE(NEW.appointment_id, -1)
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
