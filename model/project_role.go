package model

import "time"

type ProjectRole struct {
	UserID    int64     `json:"user_id" gorm:"primary_key"`
	ProjID    int64     `json:"project_id" gorm:"primary_key"`
	Role      int8      `gorm:"type:tinyint"`
	CreatedAt time.Time ``
	UpdatedAt time.Time ``
}
