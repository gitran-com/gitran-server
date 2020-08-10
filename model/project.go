package model

import (
	"time"
)

type Project struct {
	PID         uint64    `gorm:"primary_key;auto_increment"`
	Name        string    `gorm:"unique_index;size:255"`
	Type        uint8     `gorm:"type:tinyint"` //0:public 1:private
	Desc        string    `gorm:"type:text"`
	TargetLangs string    `gorm:"size:255"`
	CreatedTime time.Time ``
	UpdatedTime time.Time ``
}
