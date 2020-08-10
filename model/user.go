package model

import "time"

type User struct {
	UID            uint64    `gorm:"primary_key;auto_increment"`
	Name           string    `gorm:"unique_index;size:255"`
	Email          string    `gorm:"unique_index;size:255"`
	Password       string    `gorm:"type:char(128)"`
	Salt           string    `gorm:"type:char(128)"`
	GithubID       string    ``
	PreferredLangs string    `gorm:"type:varchar(128)"`
	CreatedTime    time.Time ``
	UpdatedTime    time.Time ``
}
