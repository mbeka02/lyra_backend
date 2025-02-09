// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package database

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users(full_name , password , email , telephone_number ,user_role) VALUES ($1,$2,$3,$4,$5) RETURNING user_id, full_name, password, email, telephone_number, created_at, user_role, verified_at, password_changed_at
`

type CreateUserParams struct {
	FullName        string
	Password        string
	Email           string
	TelephoneNumber string
	UserRole        Role
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.FullName,
		arg.Password,
		arg.Email,
		arg.TelephoneNumber,
		arg.UserRole,
	)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.FullName,
		&i.Password,
		&i.Email,
		&i.TelephoneNumber,
		&i.CreatedAt,
		&i.UserRole,
		&i.VerifiedAt,
		&i.PasswordChangedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT user_id, full_name, password, email, telephone_number, created_at, user_role, verified_at, password_changed_at FROM users WHERE email=$1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.FullName,
		&i.Password,
		&i.Email,
		&i.TelephoneNumber,
		&i.CreatedAt,
		&i.UserRole,
		&i.VerifiedAt,
		&i.PasswordChangedAt,
	)
	return i, err
}

const getUsers = `-- name: GetUsers :many
SELECT user_id , full_name , email FROM users LIMIT $1 OFFSET $2
`

type GetUsersParams struct {
	Limit  int32
	Offset int32
}

type GetUsersRow struct {
	UserID   int64
	FullName string
	Email    string
}

func (q *Queries) GetUsers(ctx context.Context, arg GetUsersParams) ([]GetUsersRow, error) {
	rows, err := q.db.QueryContext(ctx, getUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUsersRow
	for rows.Next() {
		var i GetUsersRow
		if err := rows.Scan(&i.UserID, &i.FullName, &i.Email); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users SET full_name=$1 ,email=$2 , telephone_number=$3 WHERE user_id=$4
`

type UpdateUserParams struct {
	FullName        string
	Email           string
	TelephoneNumber string
	UserID          int64
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.ExecContext(ctx, updateUser,
		arg.FullName,
		arg.Email,
		arg.TelephoneNumber,
		arg.UserID,
	)
	return err
}
