package db_test

import (
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
	"gssm/db"
	"testing"
)

type SecretGroupTestSuite struct {
	suite.Suite
	inj *do.Injector
}

func (suite *SecretGroupTestSuite) SetupSuite() {
	inj := db.CreateInjectorWithDb()
	do.Provide(inj, db.NewUserSource)
	do.Provide(inj, db.NewProjectSource)
	do.Provide(inj, db.NewSecretSource)
	suite.inj = inj
}

func (suite *SecretGroupTestSuite) TestCreateGroup() {
	t := suite.T()
	userSource := do.MustInvoke[*db.UserSource](suite.inj)
	ps := do.MustInvoke[*db.ProjectSource](suite.inj)
	ss := do.MustInvoke[*db.SecretSource](suite.inj)

	id, err := userSource.CreateUser(context.Background(), "test@test.ru", "123123123")
	if err != nil {
		return
	}

	require.NoError(t, err)

	project, err := ps.CreateProject(context.Background(), "test", id)
	require.NoError(t, err)
	require.NotEmpty(t, project)

	create, err := ss.Create(context.Background(), "test", "jwt", project)
	assert.NoError(t, err)
	assert.NotEmpty(t, create)
}

func (suite *SecretGroupTestSuite) TestGetGroup() {
	t := suite.T()
	userSource := do.MustInvoke[*db.UserSource](suite.inj)
	ps := do.MustInvoke[*db.ProjectSource](suite.inj)
	ss := do.MustInvoke[*db.SecretSource](suite.inj)

	id, err := userSource.CreateUser(context.Background(), "test@test.ru", "123123123")
	if err != nil {
		return
	}

	require.NoError(t, err)

	project, err := ps.CreateProject(context.Background(), "test", id)
	require.NoError(t, err)
	require.NotEmpty(t, project)

	create, err := ss.Create(context.Background(), "test", "jwt", project)
	require.NoError(t, err)
	require.NotEmpty(t, create)

	group, err := ss.GetById(context.Background(), create)
	require.NoError(t, err)

	assert.Equal(t, "test", group.Name)
	assert.Equal(t, "jwt", group.Prefix)
	assert.Equal(t, project, group.ProjectId)
}

func TestSecretGroupTestSuite(t *testing.T) {
	suite.Run(t, new(SecretGroupTestSuite))
}
