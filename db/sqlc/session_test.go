package db

import (
	"context"
	"testing"
	"time"

	"simplebank/util"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createRandomSession(t *testing.T, user User) Session {
	arg := CreateSessionParams{
		ID:           uuid.New(),
		Username:     user.Username,
		RefreshToken: util.RandomString(32),
		UserAgent:    util.RandomString(10),
		ClientIp:     util.RandomString(15),
		IsBlocked:    false,
		ExpiresAt:    time.Now().Add(time.Hour * 24),
	}

	session, err := testStore.CreateSession(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, session)

	require.Equal(t, arg.ID, session.ID)
	require.Equal(t, arg.Username, session.Username)
	require.Equal(t, arg.RefreshToken, session.RefreshToken)
	require.Equal(t, arg.UserAgent, session.UserAgent)
	require.Equal(t, arg.ClientIp, session.ClientIp)
	require.Equal(t, arg.IsBlocked, session.IsBlocked)
	require.WithinDuration(t, arg.ExpiresAt, session.ExpiresAt, time.Second)

	require.NotZero(t, session.CreatedAt)

	return session
}

func TestCreateSession(t *testing.T) {
	user := createRandomUser(t)
	createRandomSession(t, user)
}

func TestGetSession(t *testing.T) {
	user := createRandomUser(t)
	session1 := createRandomSession(t, user)
	session2, err := testStore.GetSession(context.Background(), session1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, session2)

	require.Equal(t, session1.ID, session2.ID)
	require.Equal(t, session1.Username, session2.Username)
	require.Equal(t, session1.RefreshToken, session2.RefreshToken)
	require.Equal(t, session1.UserAgent, session2.UserAgent)
	require.Equal(t, session1.ClientIp, session2.ClientIp)
	require.Equal(t, session1.IsBlocked, session2.IsBlocked)
	require.WithinDuration(t, session1.ExpiresAt, session2.ExpiresAt, time.Second)
	require.WithinDuration(t, session1.CreatedAt, session2.CreatedAt, time.Second)
}
