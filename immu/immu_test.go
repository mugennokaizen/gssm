package immu_test

import (
	"fmt"
	"github.com/oklog/ulid/v2"
	"github.com/samber/do"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
	"gssm/immu"
	"gssm/types"
	"testing"
)

type ImmuTestSuite struct {
	suite.Suite
	inj *do.Injector
}

func (suite *ImmuTestSuite) SetupSuite() {
	ti := immu.SetupTestImmu()
	fmt.Println(ti.Port)
	inj := do.New()
	t := suite.T()
	t.Setenv("immu.user", "immudb")
	t.Setenv("immu.password", "immudb")
	t.Setenv("immu.db", "defaultdb")
	t.Setenv("immu.port", ti.Port.Port())

	viper.AutomaticEnv()

	do.Provide(inj, immu.NewDatabase)
	do.Provide(inj, immu.NewManager)
	suite.inj = inj
}

func (suite *ImmuTestSuite) TestSet() {

	im := do.MustInvoke[*immu.Manager](suite.inj)

	t := suite.T()

	err := im.Open(context.Background())
	require.NoError(t, err)

	err = im.SetSecret(context.Background(), types.ULID(ulid.Make().String()), "group", "key", "bla-bla-bla-bla-bla-bla")
	require.NoError(t, err)

	err = im.Close(context.Background())
	require.NoError(t, err)
}

func (suite *ImmuTestSuite) TestGet() {

	val := "bla-bla-bla-bla-bla-bla"
	ul := types.ULID(ulid.Make().String())
	im := do.MustInvoke[*immu.Manager](suite.inj)

	t := suite.T()

	err := im.Open(context.Background())
	require.NoError(t, err)

	err = im.SetSecret(context.Background(), ul, "group", "key", val)
	require.NoError(t, err)

	value, err := im.GetSecret(context.Background(), ul, "group", "key")

	require.NoError(t, err)
	assert.EqualValues(t, val, value)

	err = im.Close(context.Background())
	require.NoError(t, err)
}

func TestImmuTestSuite(t *testing.T) {
	suite.Run(t, new(ImmuTestSuite))
}
