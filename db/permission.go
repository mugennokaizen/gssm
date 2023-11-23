package db

type Permission int

const (
	SecretRead   Permission = 1 << iota
	SecretModify            = 1 << iota
	SecretCreate            = 1 << iota
)
