package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/mbeka02/lyra_backend/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestGetUsers(t *testing.T) {
	users, err := store.GetUsers(context.Background(), database.GetUsersParams{
		Limit:  10,
		Offset: 0,
	})
	assert.NoError(t, err)
	fmt.Println("users=>", users)
}
