package db

import (
	"context"
	"github.com/samber/do"
	"gorm.io/gorm"
	"gssm/types"
	"gssm/utils"
)

type UserSource struct {
	db *gorm.DB
}

func NewUserSource(inj *do.Injector) (*UserSource, error) {
	db := do.MustInvoke[*gorm.DB](inj)

	return &UserSource{
		db: db,
	}, nil
}

func (us *UserSource) GetById() {

}

func (us *UserSource) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	result := us.db.WithContext(ctx).First(&user, "email = ?", email)
	return &user, result.Error
}

func (us *UserSource) CreateUser(ctx context.Context, email string, password string) (types.ULID, error) {

	hashes := utils.HashPassword(password)

	user := User{
		Email:        email,
		PasswordHash: hashes.Hash,
		Salt:         hashes.Salt,
	}

	result := us.db.WithContext(ctx).Create(&user)

	return user.Id, result.Error
}
