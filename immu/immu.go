package immu

import (
	"fmt"
	immudb "github.com/codenotary/immudb/pkg/client"
	"github.com/samber/do"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"gssm/types"
)

type Manager struct {
	client   immudb.ImmuClient
	user     []byte
	password []byte
	db       string
}

// I spent whole night to determine can I reuse session or should open it on each request. Decided to open on each request

func NewDatabase(_ *do.Injector) (immudb.ImmuClient, error) {
	host := viper.GetString("immu.host")
	port := viper.GetInt("immu.port")

	opts := immudb.DefaultOptions().
		WithAddress(host).
		WithPort(port)
	client := immudb.NewClient().WithOptions(opts)

	immuClient := immudb.ImmuClient(client)
	return immuClient, nil
}

func NewManager(inj *do.Injector) (*Manager, error) {
	client := do.MustInvoke[immudb.ImmuClient](inj)

	user := viper.GetString("immu.user")
	password := viper.GetString("immu.password")
	db := viper.GetString("immu.db")

	return &Manager{
		client:   client,
		user:     []byte(user),
		password: []byte(password),
		db:       db,
	}, nil
}

func (m *Manager) Open(ctx context.Context) error {
	return m.client.OpenSession(ctx, m.user, m.password, m.db)
}

func (m *Manager) Close(ctx context.Context) error {
	return m.client.CloseSession(ctx)
}

func (m *Manager) SetSecret(ctx context.Context, groupUlid types.ULID, group, key, value string) error {
	_, err := m.client.VerifiedSet(
		ctx,
		[]byte(fmt.Sprintf("%s.%s.%s", groupUlid, group, key)),
		[]byte(value),
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) GetSecret(ctx context.Context, groupUlid types.ULID, group, key string) (string, error) {
	entry, err := m.client.Get(
		ctx,
		[]byte(fmt.Sprintf("%s.%s.%s", groupUlid, group, key)),
	)
	if err != nil {
		return "", err
	}

	return string(entry.Value), nil
}
