package model

import (
	"time"

	"github.com/gitran-com/gitran-server/config"
	"gorm.io/gorm/clause"
)

type Role int8

const (
	//RoleAdmin can view, submit, vote and commit translations. Can also manage project settings, including adding collaborators
	RoleAdmin Role = iota
	//RoleCommitter can view, submit, vote and commit translations
	RoleCommitter
	//RoleContributor can view, submit and vote translations
	RoleContributor
	//RoleViewer can view translations
	RoleViewer
)

type ProjRole struct {
	ProjID    int64     `json:"proj_id" gorm:"primaryKey"`
	UserID    int64     `json:"user_id" gorm:"primaryKey"`
	Role      Role      `json:"role" gorm:"type:tinyint"`
	CreatedAt time.Time ``
	UpdatedAt time.Time ``
}

//TableName return table name
func (*ProjRole) TableName() string {
	return config.DB.TablePrefix + "project_roles"
}

//Write writes project to DB
func (role *ProjRole) Write() {
	db.Save(role)
}

//Create creates new project
func (role *ProjRole) Create() error {
	return db.Create(role).Error
}

func GetUserProjRole(user_id int64, proj_id int64) *ProjRole {
	var roles []ProjRole
	db.Where("proj_id=? AND user_id=?", proj_id, user_id).First(&roles)
	if len(roles) > 0 {
		return &roles[0]
	}
	return nil
}

func SetUserProjRole(user_id int64, proj_id int64, role Role) {
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "proj_id"}}, // key colume
		DoUpdates: clause.AssignmentColumns([]string{"role"}),            // column needed to be updated
	}).Create(&ProjRole{
		UserID: user_id,
		ProjID: proj_id,
		Role:   role,
	})
}
