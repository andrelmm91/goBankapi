package db

import (
	"context"
	"simplebank/util"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func createRandomVerifyEmail(t *testing.T) VerifyEmail {
	user1 := createRandomUser(t)

	arg := CreateVerifyEmailParams{
		Username:   user1.Username,
		Email:      user1.Email,
		SecretCode: util.RandomString(32),
	}

	verifyEmail, err := testStore.CreateVerifyEmail(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, verifyEmail)

	require.Equal(t, arg.Username, verifyEmail.Username)
	require.Equal(t, arg.Email, verifyEmail.Email)
	require.Equal(t, arg.SecretCode, verifyEmail.SecretCode)
	require.False(t, verifyEmail.IsUsed)
	require.WithinDuration(t, time.Now(), verifyEmail.CreatedAt, time.Second)
	require.WithinDuration(t, time.Now().Add(15*time.Minute), verifyEmail.ExpiredAt, time.Second)

	return verifyEmail
}

func TestCreateVerifyEmail(t *testing.T) {
	createRandomVerifyEmail(t)
}

func TestUpdateVerifyEmail(t *testing.T) {
	verifyEmail := createRandomVerifyEmail(t)

	arg := UpdateVerifyEmailParams{
		ID:         verifyEmail.ID,
		SecretCode: verifyEmail.SecretCode,
	}

	updatedVerifyEmail, err := testStore.UpdateVerifyEmail(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedVerifyEmail)

	require.Equal(t, verifyEmail.ID, updatedVerifyEmail.ID)
	require.Equal(t, verifyEmail.Username, updatedVerifyEmail.Username)
	require.Equal(t, verifyEmail.Email, updatedVerifyEmail.Email)
	require.Equal(t, verifyEmail.SecretCode, updatedVerifyEmail.SecretCode)
	require.True(t, updatedVerifyEmail.IsUsed)
	require.WithinDuration(t, verifyEmail.CreatedAt, updatedVerifyEmail.CreatedAt, time.Second)
	require.WithinDuration(t, verifyEmail.ExpiredAt, updatedVerifyEmail.ExpiredAt, time.Second)
}
