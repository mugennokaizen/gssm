package db_test

import (
	"context"
	"github.com/samber/do"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gssm/db"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testDbInstance *gorm.DB
var injector *do.Injector

func TestMain(m *testing.M) {
	testDB := db.SetupTestDatabase()

	testDbInstance = testDB.DbInstance

	injector = do.New()
	do.ProvideValue(injector, testDbInstance)
	do.Provide(injector, db.NewUserSource)

	os.Exit(m.Run())
}

func TestCreateUser(t *testing.T) {
	userSource := do.MustInvoke[*db.UserSource](injector)

	_, err := userSource.CreateUser(context.Background(), "test@test.ru", "123123123")
	if err != nil {
		return
	}

	require.NoError(t, err)
}

func TestUserExist(t *testing.T) {
	userSource := do.MustInvoke[*db.UserSource](injector)
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

func TestFindUser(t *testing.T) {
	userSource := do.MustInvoke[*db.UserSource](injector)

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
	assert.Equal(t, email.Email, "test@test.ru")
	assert.Equal(t, len(email.Salt), 8)
	assert.Equal(t, len(email.PasswordHash), 64)
}
