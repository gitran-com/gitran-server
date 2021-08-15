package model

import "github.com/gitran-com/gitran-server/config"

type ProjAccess struct {
	ID               int64 `json:"id" gorm:"primaryKey"`
	PublicView       bool  `json:"public_view"`
	PublicContribute bool  `json:"public_contribute"`
}

//TableName return table name
func (*ProjAccess) TableName() string {
	return config.DB.TablePrefix + "project_access"
}

//Write writes project to DB
func (pa *ProjAccess) Write() {
	db.Save(pa)
}

//Create creates new project
func (pa *ProjAccess) Create() error {
	return db.Create(pa).Error
}

//NewProjAccess creates a new access for project
func NewProjAccess(pa *ProjAccess) (*ProjAccess, error) {
	if result := db.Create(pa); result.Error != nil {
		return nil, result.Error
	}
	return pa, nil
}
