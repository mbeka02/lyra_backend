package repository

import (
	"context"
	"testing"
	"time"

	"github.com/mbeka02/lyra_backend/internal/database"
	"github.com/mbeka02/lyra_backend/internal/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) database.User {
	params := database.CreateUserParams{
		FullName:        util.RandName(),
		TelephoneNumber: util.RandPhoneNumber(),
		UserRole:        database.Role(util.RandRole()),
		DateOfBirth:     util.RandDateOfBirth(16, 65),
		Email:           util.RandEmail(),
		Password:        "mySuperDuperSecretPassword",
	}
	user, err := store.CreateUser(context.Background(), params)
	require.NoError(t, err)

	require.NotEmpty(t, user)
	require.Equal(t, params.FullName, user.FullName)
	require.Equal(t, params.Email, user.Email)
	require.Equal(t, params.TelephoneNumber, user.TelephoneNumber)
	require.Equal(t, params.DateOfBirth, user.DateOfBirth)
	require.Equal(t, params.Password, user.Password)
	require.Equal(t, params.UserRole, user.UserRole)
	require.NotZero(t, user.CreatedAt)
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUserByEmail(t *testing.T) {
	randomUser := createRandomUser(t)
	user, err := store.GetUserByEmail(context.Background(), randomUser.Email)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, randomUser.FullName, user.FullName)
	require.Equal(t, randomUser.Password, user.Password)
	require.Equal(t, randomUser.Email, user.Email)
	require.Equal(t, randomUser.TelephoneNumber, user.TelephoneNumber)
	require.Equal(t, randomUser.DateOfBirth, user.DateOfBirth)
	require.Equal(t, randomUser.UserRole, user.UserRole)

	require.WithinDuration(t, randomUser.CreatedAt, user.CreatedAt, time.Second)
}

func TestGetUserById(t *testing.T) {
	randomUser := createRandomUser(t)
	user, err := store.GetUserById(context.Background(), randomUser.UserID)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, randomUser.FullName, user.FullName)
	require.Equal(t, randomUser.Password, user.Password)
	require.Equal(t, randomUser.Email, user.Email)
	require.Equal(t, randomUser.TelephoneNumber, user.TelephoneNumber)
	require.Equal(t, randomUser.DateOfBirth, user.DateOfBirth)
	require.Equal(t, randomUser.UserRole, user.UserRole)

	require.WithinDuration(t, randomUser.CreatedAt, user.CreatedAt, time.Second)
}
