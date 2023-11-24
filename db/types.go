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

type Project struct {
	Identity
	Name      string     `json:"name"`
	CreatorId types.ULID `json:"-"`
}

type UserToProject struct {
	Identity
	ProjectId  types.ULID
	UserId     types.ULID
	Permission Permission
}

type SecretGroup struct {
	Identity
	Name      string
	Prefix    string
	ProjectId types.ULID
}
