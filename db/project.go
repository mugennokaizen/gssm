package db

import (
	"context"
	"github.com/samber/do"
	"gorm.io/gorm"
	"gssm/types"
)

type ProjectSource struct {
	db *gorm.DB
}

func NewProjectSource(inj *do.Injector) (*ProjectSource, error) {
	db := do.MustInvoke[*gorm.DB](inj)

	return &ProjectSource{
		db: db,
	}, nil
}

func (ps *ProjectSource) CreateProject(ctx context.Context, name string, userId types.ULID) (types.ULID, error) {
	newProject := Project{
		Name:      name,
		CreatorId: userId,
	}

	err := ps.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(&newProject).Error; err != nil {
			return err
		}

		if err := tx.Create(&UserToProject{
			Identity:   Identity{},
			ProjectId:  newProject.Id,
			UserId:     userId,
			Permission: SecretRead | SecretCreate | SecretModify,
		}).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	return newProject.Id, nil
}

func (ps *ProjectSource) GetProjects(ctx context.Context, userId types.ULID) ([]*Project, error) {
	var projects []*Project

	ps.db.
		WithContext(ctx).
		Table("project as p").
		Select("p.id, p.name, p.creator_id").
		Joins("inner join user_to_project as ut on ut.project_id = p.id").
		Where("ut.user_id = ?", userId).Find(&projects)

	return projects, nil
}

func (ps *ProjectSource) GetProjectPermissions(ctx context.Context, projectId types.ULID, userId types.ULID) Permission {
	var permission Permission

	ps.db.
		WithContext(ctx).
		Table("user_to_project as ut").
		Select("permission").
		Where("ut.user_id = ? and ut.project_id = ?", userId, projectId).Scan(&permission)

	return permission
}
