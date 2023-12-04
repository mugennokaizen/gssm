package db_test

import (
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
	"gssm/db"
	"gssm/utils"
	"testing"
)

type AccessKeyTestSuite struct {
	suite.Suite
	inj *do.Injector
	ctx context.Context
}

func (suite *AccessKeyTestSuite) SetupSuite() {
	inj := db.CreateInjectorWithDb()
	do.Provide(inj, db.NewUserSource)
	do.Provide(inj, db.NewProjectSource)
	do.Provide(inj, db.NewAccessKeySource)
	suite.inj = inj
	suite.ctx = context.Background()
}

func (suite *AccessKeyTestSuite) TestCreateAccessKey() {
	t := suite.T()
	userSource := do.MustInvoke[*db.UserSource](suite.inj)
	ps := do.MustInvoke[*db.ProjectSource](suite.inj)
	aks := do.MustInvoke[*db.AccessKeySource](suite.inj)

	id, err := userSource.CreateUser(suite.ctx, "test@test.ru", "123123123")
	if err != nil {
		return
	}

	require.NoError(t, err)

	project, err := ps.CreateProject(suite.ctx, "test", id)
	require.NoError(t, err)
	require.NotEmpty(t, project)

	key := utils.GenerateAccessKey()
	secret, mask := utils.GenerateSecretKey()

	signature, err := utils.GetSignature(id, project, secret)
	require.NoError(t, err)

	accessKey, err := aks.Create(suite.ctx, key, mask, id, project, signature, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, accessKey)
}

func (suite *AccessKeyTestSuite) TestGetAccessKey() {
	t := suite.T()

	userSource := do.MustInvoke[*db.UserSource](suite.inj)
	ps := do.MustInvoke[*db.ProjectSource](suite.inj)
	aks := do.MustInvoke[*db.AccessKeySource](suite.inj)

	id, err := userSource.CreateUser(suite.ctx, "test@test.ru", "123123123")
	if err != nil {
		return
	}

	require.NoError(t, err)

	project, err := ps.CreateProject(suite.ctx, "test", id)
	require.NoError(t, err)
	require.NotEmpty(t, project)

	key := utils.GenerateAccessKey()
	secret, mask := utils.GenerateSecretKey()

	signature, err := utils.GetSignature(id, project, secret)
	require.NoError(t, err)

	accessKey, err := aks.Create(suite.ctx, key, mask, id, project, signature, nil)
	require.NoError(t, err)
	require.NotEmpty(t, accessKey)

	accessKeyDb, err := aks.Get(suite.ctx, key)
	require.NoError(t, err)

	assert.Equal(t, key, accessKeyDb.Key)
	assert.EqualValues(t, project, accessKeyDb.ProjectId)
	assert.EqualValues(t, id, accessKeyDb.UserId)
	assert.Equal(t, mask, accessKeyDb.Mask)
	assert.Nil(t, accessKeyDb.Expires)
	assert.Len(t, accessKeyDb.Signature, 64)
}

func TestAccessKeyTestSuite(t *testing.T) {
	suite.Run(t, new(AccessKeyTestSuite))
}
