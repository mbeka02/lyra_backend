-- name: GetUsers :many
SELECT user_id , full_name , email FROM users LIMIT $1 OFFSET $2;
-- name: GetUserByEmail :one
SELECT * FROM users WHERE email=$1;
-- name: CreateUser :one
INSERT INTO users(full_name , password , email , telephone_number ,user_role) VALUES ($1,$2,$3,$4,$5) RETURNING *;
-- name: UpdateUser :one
UPDATE users SET full_name=$1 ,email=$2 , telephone_number=$3 WHERE user_id=$4 RETURNING *; 
