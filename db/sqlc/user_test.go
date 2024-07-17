package db

import (
	"context"
	"database/sql"
	"simplebank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassord, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:    util.RandomOwner(),
		HashedPassword:  hashedPassord,
		FullName: util.RandomOwner(),
		Email: util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	// testing using package Testify
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	// testing using package Testify
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUser(t *testing.T) {
	user1 := createRandomUser(t)
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := UpdateUserParams{
		HashedPassword: sql.NullString{
			String: hashedPassword,
			Valid:  true,
		},
		PasswordChangedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		FullName: sql.NullString{
			String: util.RandomOwner(),
			Valid:  true,
		},
		Email: sql.NullString{
			String: util.RandomEmail(),
			Valid:  true,
		},
		IsEmailVerified: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
		Username: sql.NullString{
			String: user1.Username,
			Valid:  true,
		},
	}

	user2, err := testQueries.UpdateUser(context.Background(), arg)

	// testing using package Testify
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, hashedPassword, user2.HashedPassword)
	require.Equal(t, arg.FullName.String, user2.FullName)
	require.Equal(t, arg.Email.String, user2.Email)
	require.Equal(t, arg.IsEmailVerified.Bool, user2.IsEmailVerified)
	require.WithinDuration(t, arg.PasswordChangedAt.Time, user2.PasswordChangedAt, time.Second)
}
