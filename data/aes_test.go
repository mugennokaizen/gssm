package data_test

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gssm/data"
	"testing"
)

type AesTestSuite struct {
	suite.Suite
	ps *data.AesProcessor
}

func (suite *AesTestSuite) SetupTest() {
	key, _ := hex.DecodeString("4ebca8b93b2fa70067ea92b33cdf9b4d011e3acb2ccabe768272157058e1cb9b")

	suite.ps = &data.AesProcessor{
		MasterKey: key,
	}
}

func (suite *AesTestSuite) TestEncryptDecrypt() {
	t := suite.T()

	value := "host=localhost user=user password=password dbname=db port=port sslmode=disable"

	encryptedData, err := suite.ps.Encrypt(value)
	require.NoError(t, err)

	decrypt, err := suite.ps.Decrypt(encryptedData)
	require.NoError(t, err)

	assert.Equal(t, value, decrypt)
}

func TestAesTestSuite(t *testing.T) {
	suite.Run(t, new(AesTestSuite))
}
