package model

type File struct {
	FileID uint64 `gorm:"primary_key;auto_increment"`
	ProjID uint64 `json:"project_id"`
	Dir    string `gorm:"type:varchar(256);not null"`
	Name   string `gorm:"type:varchar(256);not null"`
}
