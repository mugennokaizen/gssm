package db_test

import (
	"context"
	"github.com/samber/do"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gssm/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UserTestSuite struct {
	suite.Suite
	inj *do.Injector
}

func (suite *UserTestSuite) SetupSuite() {
	inj := db.CreateInjectorWithDb()
	do.Provide(inj, db.NewUserSource)
	suite.inj = inj
}

func (suite *UserTestSuite) TestCreateUser() {
	t := suite.T()
	userSource := do.MustInvoke[*db.UserSource](suite.inj)

	_, err := userSource.CreateUser(context.Background(), "test@test.ru", "123123123")
	if err != nil {
		return
	}

	require.NoError(t, err)
}

func (suite *UserTestSuite) TestUserExist() {
	t := suite.T()
	userSource := do.MustInvoke[*db.UserSource](suite.inj)
	email := "test@test.ru"
	_, err := userSource.CreateUser(context.Background(), email, "123123123")
	if err != nil {
		return
	}

	require.NoError(t, err)

	res := userSource.IsUserExist(context.Background(), email)
	assert.True(t, res)

	res2 := userSource.IsUserExist(context.Background(), "adasdad@test.com")

	assert.False(t, res2)
}

func (suite *UserTestSuite) TestFindUser() {
	t := suite.T()
	userSource := do.MustInvoke[*db.UserSource](suite.inj)

	_, err := userSource.CreateUser(context.Background(), "test@test.ru", "123123123")
	if err != nil {
		return
	}

	require.NoError(t, err)

	email, err := userSource.GetUserByEmail(context.Background(), "test@test.ru")
	if err != nil {
		return
	}
	require.NoError(t, err)
	assert.Equal(t, "test@test.ru", email.Email)
	assert.Equal(t, 8, len(email.Salt))
	assert.Equal(t, 64, len(email.PasswordHash))
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
