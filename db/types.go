package db

import "gssm/types"

type Identity struct {
	Id types.ULID `json:"-" gorm:"primaryKey;type:ulid;default:gen_ulid()"`
}

type User struct {
	Identity
	Email        string `json:"email"`
	PasswordHash []byte `json:"-"`
	Salt         []byte `json:"-"`
}
