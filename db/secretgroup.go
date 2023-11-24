package db

import (
	"github.com/samber/do"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"gssm/types"
)

type SecretSource struct {
	db *gorm.DB
}

func NewSecretSource(inj *do.Injector) (*SecretSource, error) {
	db := do.MustInvoke[*gorm.DB](inj)

	return &SecretSource{
		db: db,
	}, nil
}

func (ps *SecretSource) Create(ctx context.Context, name, prefix string, projectId types.ULID) (types.ULID, error) {
	group := SecretGroup{
		Name:      name,
		Prefix:    prefix,
		ProjectId: projectId,
	}

	result := ps.db.WithContext(ctx).Create(&group)
	return group.Id, result.Error
}

func (ps *SecretSource) GetById(ctx context.Context, id types.ULID) (*SecretGroup, error) {
	var group SecretGroup
	result := ps.db.WithContext(ctx).First(&group, "id = ?", id)
	return &group, result.Error
}
