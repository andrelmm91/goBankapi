package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestVerifyEmailTx tests the VerifyEmailTx function
func TestVerifyEmailTx(t *testing.T) {
	store := NewStore(testDB)

	verifyEmail := createRandomVerifyEmail(t)
	arg := VerifyEmailTxParams{
		EmailId:    verifyEmail.ID,
		SecretCode: verifyEmail.SecretCode,
	}

	result, err := store.VerifyEmailTx(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	// Verify VerifyEmail fields
	require.Equal(t, verifyEmail.ID, result.VerifyEmail.ID)
	require.Equal(t, verifyEmail.Username, result.VerifyEmail.Username)
	require.Equal(t, verifyEmail.Email, result.VerifyEmail.Email)
	require.Equal(t, verifyEmail.SecretCode, result.VerifyEmail.SecretCode)
	require.True(t, result.VerifyEmail.IsUsed)
	require.WithinDuration(t, verifyEmail.CreatedAt, result.VerifyEmail.CreatedAt, time.Second)
	require.WithinDuration(t, verifyEmail.ExpiredAt, result.VerifyEmail.ExpiredAt, time.Second)

	// Verify User fields
	require.Equal(t, verifyEmail.Username, result.User.Username)
	require.True(t, result.User.IsEmailVerified)
}
