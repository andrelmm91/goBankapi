package db

import (
	"context"
	"simplebank/util"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

// createRandomUserParams creates random parameters for creating a user
func createRandomUserParams(t *testing.T) CreateUserParams {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	return CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
}

// TestCreateUserTx tests the CreateUserTx function
func TestCreateUserTx(t *testing.T) {
	store := NewStore(testDB)

	arg := CreateUserTxParams{
		CreateUserParams: createRandomUserParams(t),
		AfterCreate: func(user User) error {
			// Placeholder for additional operations after user creation
			// For now, we just return nil to indicate success
			return nil
		},
	}

	ctx := context.Background()
	result, err := store.CreateUserTx(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	// Verify User fields
	user := result.User
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.WithinDuration(t, time.Now(), user.CreatedAt, time.Second)
	require.False(t, user.IsEmailVerified)
}
