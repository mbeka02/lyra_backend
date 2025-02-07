package auth

import (
	"testing"

	"github.com/mbeka02/lyra_backend/internal/util"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordHashing(t *testing.T) {
	password := util.RandString(10)
	hashedPassword, err := HashPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	err = ComparePassword(password, hashedPassword)
	require.NoError(t, err)

	wrongPassword := util.RandString(10)
	err = ComparePassword(wrongPassword, hashedPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
