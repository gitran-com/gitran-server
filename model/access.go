package model

import (
	"time"

	"github.com/gitran-com/gitran-server/constant"
)

type ProjectRole struct {
	UserID    int64         `json:"user_id" gorm:"primary_key"`
	ProjID    int64         `json:"project_id" gorm:"primary_key"`
	Role      constant.Role `json:"role" gorm:"type:tinyint"`
	CreatedAt time.Time     ``
	UpdatedAt time.Time     ``
}
