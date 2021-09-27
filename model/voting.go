package model

import (
	"time"

	"github.com/gitran-com/gitran-server/config"
)

type Voting struct {
	UserID    int64     `gorm:"primaryKey"`
	TranID    int64     `gorm:"primaryKey"`
	Vote      int       `gorm:"vote"`
	CreatedAt time.Time ``
	UpdatedAt time.Time ``
	DeletedAt time.Time ``
}

//TableName returns table name
func (*Voting) TableName() string {
	return config.DB.TablePrefix + "votings"
}
