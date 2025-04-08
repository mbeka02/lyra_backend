-- name: GetUsers :many
SELECT user_id , full_name , email FROM users LIMIT $1 OFFSET $2;
-- name: GetUserByEmail :one
SELECT * FROM users WHERE email=$1;
-- name: GetUserById :one
SELECT * FROM users WHERE user_id=$1;
-- name: CreateUser :one
INSERT INTO users(full_name , password ,date_of_birth, email , telephone_number ,user_role) VALUES ($1,$2,$3,$4,$5,$6) RETURNING *;
-- name: UpdateUser :exec
UPDATE users SET full_name=$1 ,email=$2 , telephone_number=$3 WHERE user_id=$4;

-- name: UpdateProfilePicture :exec
UPDATE users SET profile_image_url=$1 WHERE user_id=$2;

-- name: CompleteOnboarding :exec
UPDATE users SET is_onboarded=true WHERE user_id=$1;
