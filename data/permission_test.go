package data_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gssm/data"
	"gssm/db"
	"testing"
)

type PermissionTestSuite struct {
	suite.Suite
	ps *data.PermissionProcessor
}

func (suite *PermissionTestSuite) SetupTest() {
	suite.ps = &data.PermissionProcessor{}
}

func (suite *PermissionTestSuite) TestVerifyPermission() {
	t := suite.T()

	p1 := db.SecretCreate | db.SecretModify | db.SecretRead

	assert.True(t, suite.ps.CheckPermission(p1, db.SecretCreate))
	assert.True(t, suite.ps.CheckPermission(p1, db.SecretRead))
	assert.True(t, suite.ps.CheckPermission(p1, db.SecretModify))

	p2 := db.SecretRead

	assert.True(t, suite.ps.CheckPermission(p2, db.SecretRead))
	assert.False(t, suite.ps.CheckPermission(p2, db.SecretCreate))
}

func (suite *PermissionTestSuite) TestRemovePermission() {
	t := suite.T()

	p1 := db.SecretCreate | db.SecretModify | db.SecretRead

	assert.EqualValues(t, suite.ps.RemovePermission(p1, db.SecretCreate), db.SecretModify|db.SecretRead)
	assert.EqualValues(t, suite.ps.RemovePermission(p1, db.SecretModify), db.SecretCreate|db.SecretRead)

	p2 := db.SecretRead

	assert.EqualValues(t, suite.ps.RemovePermission(p2, db.SecretRead), 0)
}

func (suite *PermissionTestSuite) TestAddPermission() {
	t := suite.T()

	var p1 db.Permission = db.SecretCreate

	assert.EqualValues(t, suite.ps.AddPermission(p1, db.SecretRead), db.SecretCreate|db.SecretRead)
	assert.EqualValues(t, suite.ps.AddPermission(p1, db.SecretCreate), db.SecretCreate)
	assert.EqualValues(t, suite.ps.AddPermission(p1, db.SecretModify), db.SecretCreate|db.SecretModify)

}

func TestPermissionTestSuite(t *testing.T) {
	suite.Run(t, new(PermissionTestSuite))
}
