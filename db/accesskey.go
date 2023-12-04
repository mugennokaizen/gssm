package db

import (
	"github.com/samber/do"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"gssm/types"
	"time"
)

type AccessKeySource struct {
	db *gorm.DB
}

func NewAccessKeySource(inj *do.Injector) (*AccessKeySource, error) {
	db := do.MustInvoke[*gorm.DB](inj)

	return &AccessKeySource{
		db: db,
	}, nil
}

func (ps *AccessKeySource) Create(ctx context.Context, key, mask string, userId, projectId types.ULID, signature []byte, expires *time.Time) (types.ULID, error) {
	group := AccessKey{
		Identity:  Identity{},
		ProjectId: projectId,
		UserId:    userId,
		Mask:      mask,
		Key:       key,
		Signature: signature,
		Expires:   expires,
	}

	result := ps.db.WithContext(ctx).Create(&group)
	return group.Id, result.Error
}

func (ps *AccessKeySource) Get(ctx context.Context, key string) (*AccessKey, error) {
	var group AccessKey
	result := ps.db.WithContext(ctx).First(&group, "key = ?", key)

	return &group, result.Error
}
