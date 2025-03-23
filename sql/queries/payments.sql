-- name: CreatePayment :one
INSERT INTO payments (
  reference,
  amount,
  current_status,
  created_at,
  patient_id,
  doctor_id,
  appointment_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: UpdatePaymentStatus :exec
UPDATE payments
SET 
  current_status = $1,
  completed_at = CASE WHEN $1 = 'completed' THEN NOW() ELSE completed_at END,
  updated_at = NOW()
WHERE reference = $2;

-- name: GetPaymentByReference :one
SELECT * FROM payments WHERE reference = $1 LIMIT 1;
