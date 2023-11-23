package data

import (
	"github.com/samber/do"
	"gssm/db"
)

type PermissionProcessor struct{}

func NewPermissionProcessor(_ *do.Injector) (*PermissionProcessor, error) {
	return &PermissionProcessor{}, nil
}

func (t *PermissionProcessor) CheckPermission(userPermissions db.Permission, p db.Permission) bool {
	return userPermissions&p == p
}

func (t *PermissionProcessor) RemovePermission(userPermissions db.Permission, p db.Permission) db.Permission {
	return userPermissions ^ p
}

func (t *PermissionProcessor) AddPermission(userPermission db.Permission, p db.Permission) db.Permission {
	return userPermission | p
}
