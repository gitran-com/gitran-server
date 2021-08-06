package model

import "time"

type Authority struct {
	UserID    int64     `json:"user_id" gorm:"primary_key;auto_increment"`
	ProjID    int64     `json:"project_id" gorm:"primary_key;auto_increment"`
	Role      uint8     `gorm:"type:tinyint"`
	CreatedAt time.Time ``
	UpdatedAt time.Time ``
}
