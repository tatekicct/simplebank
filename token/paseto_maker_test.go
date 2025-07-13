package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tatekicct/simplebank/util"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	role := util.DepositorRole
	duration := time.Minute

	token, err := maker.CreateToken(username, role, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomOwner(), util.DepositorRole, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestPasetoWrongTokenType(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomOwner(), util.DepositorRole, time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// payload, err := maker.VerifyToken(token)
	// require.Error(t, err)
	// require.EqualError(t, err, ErrInvalidToken.Error())
	// require.Nil(t, payload)
}
