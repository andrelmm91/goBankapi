package db

import (
	"context"
	"simplebank/util"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassord, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassord,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testStore.CreateUser(context.Background(), arg)

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
	user2, err := testStore.GetUser(context.Background(), user1.Username)

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
		HashedPassword: pgtype.Text{
			String: hashedPassword,
			Valid:  true,
		},
		PasswordChangedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
		FullName: pgtype.Text{
			String: util.RandomOwner(),
			Valid:  true,
		},
		Role: pgtype.Text{
			String: util.DepositorRole,
			Valid:  true,
		},
		Email: pgtype.Text{
			String: util.RandomEmail(),
			Valid:  true,
		},
		IsEmailVerified: pgtype.Bool{
			Bool:  true,
			Valid: true,
		},
		Username: user1.Username,
	}

	user2, err := testStore.UpdateUser(context.Background(), arg)

	// testing using package Testify
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, hashedPassword, user2.HashedPassword)
	require.Equal(t, arg.FullName.String, user2.FullName)
	require.Equal(t, arg.Email.String, user2.Email)
	require.Equal(t, arg.Role.String, user2.Role)
	require.Equal(t, arg.IsEmailVerified.Bool, user2.IsEmailVerified)
	require.WithinDuration(t, arg.PasswordChangedAt.Time, user2.PasswordChangedAt, time.Second)
}
