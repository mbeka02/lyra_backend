-- name: CreatePayment :one
INSERT INTO payments (
  reference,
  amount,
  patient_id,
  doctor_id,
  appointment_id
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: UpdatePaymentStatus :exec
UPDATE payments
SET 
  current_status = $1,
  --NB: Cast string literal to the appropriate type (payment status)
  completed_at = CASE WHEN $1 = 'completed'::payment_status THEN NOW() ELSE completed_at END,
  updated_at = NOW()
WHERE reference = $2;

-- name: GetPaymentByReference :one
SELECT * FROM payments WHERE reference = $1 LIMIT 1;
