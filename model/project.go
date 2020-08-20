package model

import (
	"time"
)

type Project struct {
	ID          uint64    `gorm:"primary_key;auto_increment"`
	Name        string    `gorm:"unique_index;size:255"`
	Desc        string    `gorm:"type:text"`
	Type        uint8     `gorm:"type:tinyint"` //0:public 1:private
	TargetLangs string    `gorm:"size:255"`
	CreatedAt   time.Time ``
	UpdatedAt   time.Time ``
}
