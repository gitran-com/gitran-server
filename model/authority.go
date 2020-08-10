package model

import "time"

type Authority struct {
	UID         uint64    `gorm:"primary_key;auto_increment"`
	PID         uint64    `gorm:"primary_key;auto_increment"`
	Type        uint8     `gorm:"type:tinyint"`
	CreatedTime time.Time ``
	UpdatedTime time.Time ``
}
