package token

import (
	"simplebank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	role := util.DepositorRole
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(username, role, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.Equal(t, role, payload.Role)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestPasetoKeyLengthMaker(t *testing.T) {
	// KeyLength should be equal to 32
	_, err := NewPasetoMaker(util.RandomString(31))
	require.Error(t, err)
}

func TestExpiredInvalidPasetoToken(t *testing.T) {
	role := util.DepositorRole
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(util.RandomOwner(), role, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	// Expired
	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)

	// Invalid
	tokenInvalid := util.RandomString(32)
	payload, err = maker.VerifyToken(tokenInvalid)
	require.Nil(t, payload)
	require.EqualError(t, err, ErrInvalidToken.Error())
}
