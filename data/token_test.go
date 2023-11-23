package data_test

import (
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gssm/data"
	"gssm/types"
	"os"
	"testing"
	"time"
)

var tp *data.TokenProcessor
var ul types.ULID

func TestMain(m *testing.M) {
	ul = types.ULID(ulid.Make().String())
	tp = &data.TokenProcessor{
		RefreshCookieName:    "refresh_cookie",
		RefreshTokenDuration: time.Minute * 60,
		AccessCookieName:     "access_cookie",
		AccessTokenDuration:  time.Minute * 30,
		IsSecure:             false,
		SecretKey:            "fZ_hn8W6P5sQF_Tq3MziOvhdm2hLxopDqHWUZKIe8vXngUzYUkcfQbXpJJ31U_-zezIcUBrRY-k-cY3gR2m7sZFyRXutSZJMUNfTaYWYCcB-Pw0Mo-NCGSi7UqD41alV9YeH5zWifCBLNCLQftsNC6ZXaCaT4sAHsRqrp8JxTZpFMaKNniSMaTIA1f9-AEkokkPT5awTRFAfQHte-6WquCydl2NGEFEvgONawsGhIzG6Dn5W0Fm6pnqGEch-O4JCsIRNJufTllgXdF1-fjwDoduTkp3AihPvD-Y3i6LUcBjEs2k3ckKcT41jHW9FdT5wlSqFSE009p4UUJcjbrXTxA",
	}

	os.Exit(m.Run())
}

func TestGenerateToken(t *testing.T) {
	_, err := tp.GenerateToken(types.ULID(ul), tp.RefreshTokenDuration)
	if err != nil {
		return
	}
	require.NoError(t, err)
}

func TestVerifyToken(t *testing.T) {
	token, err := tp.GenerateToken(types.ULID(ul), tp.RefreshTokenDuration)
	if err != nil {
		return
	}
	require.NoError(t, err)

	verifyToken, t2, err := tp.VerifyToken(token)

	require.NoError(t, err)
	assert.True(t, verifyToken)
	assert.Equal(t, ul, t2.Id)
}
