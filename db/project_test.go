package db_test

import (
	"context"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gssm/db"
	"testing"
)

type ProjectTestSuite struct {
	suite.Suite
	inj *do.Injector
}

func (suite *ProjectTestSuite) SetupSuite() {
	inj := db.CreateInjectorWithDb()
	do.Provide(inj, db.NewUserSource)
	do.Provide(inj, db.NewProjectSource)
	suite.inj = inj
}

func (suite *ProjectTestSuite) TestCreateProject() {
	t := suite.T()
	userSource := do.MustInvoke[*db.UserSource](suite.inj)
	ps := do.MustInvoke[*db.ProjectSource](suite.inj)

	id, err := userSource.CreateUser(context.Background(), "test@test.ru", "123123123")
	if err != nil {
		return
	}

	require.NoError(t, err)

	project, err := ps.CreateProject(context.Background(), "test", id)
	require.NoError(t, err)
	require.NotEmpty(t, project)
}

func (suite *ProjectTestSuite) TestCreatedPermissionsInProject() {
	t := suite.T()
	userSource := do.MustInvoke[*db.UserSource](suite.inj)
	ps := do.MustInvoke[*db.ProjectSource](suite.inj)

	id, err := userSource.CreateUser(context.Background(), "test@test.ru", "123123123")
	if err != nil {
		return
	}

	require.NoError(t, err)

	_, err = ps.CreateProject(context.Background(), "test", id)
	require.NoError(t, err)

	projects := ps.GetProjects(context.Background(), id)

	require.Len(t, projects, 1)

	for _, project := range projects {
		permission := ps.GetProjectPermissions(context.Background(), project.Id, id)
		assert.EqualValues(t, db.SecretCreate|db.SecretModify|db.SecretRead, permission)
	}

}

func (suite *ProjectTestSuite) TestChangeProjectName() {
	t := suite.T()
	userSource := do.MustInvoke[*db.UserSource](suite.inj)
	ps := do.MustInvoke[*db.ProjectSource](suite.inj)

	id, err := userSource.CreateUser(context.Background(), "test@test.ru", "123123123")
	if err != nil {
		return
	}

	require.NoError(t, err)

	_, err = ps.CreateProject(context.Background(), "test", id)
	require.NoError(t, err)

	projects := ps.GetProjects(context.Background(), id)

	require.Len(t, projects, 1)

	for _, project := range projects {
		err = ps.ChangeProjectName(context.Background(), project.Id, id, "blabla")

		require.NoError(t, err)
	}

	projects = ps.GetProjects(context.Background(), id)

	for _, project := range projects {
		assert.Equal(t, "blabla", project.Name)
	}
}

func TestProjectTestSuite(t *testing.T) {
	suite.Run(t, new(ProjectTestSuite))
}
