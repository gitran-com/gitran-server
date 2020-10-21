package model

import (
	"regexp"
	"time"
)

type Project struct {
	ProjID    uint64    `json:"project_id" gorm:"primary_key;auto_increment"`
	Name      string    `gorm:"type:varchar(64);unique_index;not null"`
	Desc      string    `json:"description" gorm:"type:varchar(256)"`
	Private   bool      ``
	Type      string    `gorm:"type:char(3);not null"`
	Creator   uint64    `gorm:"not null"`
	SrcLangs  string    `json:"source_languages" gorm:"type:varchar(128);not null"`
	TgtLangs  string    `json:"target_languages" gorm:"type:varchar(128)"`
	Path      string    `gorm:"type:varchar(256)"`
	CreatedAt time.Time ``
	UpdatedAt time.Time ``
}

var (
	projReg = "^[A-Za-z0-9-]{1,64}$"
)

func CheckProjName(name string) bool {
	ok, _ := regexp.Match(projReg, []byte(name))
	return ok
}

func GetProjByName(name *string) *Project {
	proj := &Project{}
	DB.Where("name=?", *name).First(proj)
	return proj
}

func NewProject(proj *Project) (*Project, error) {
	//fmt.Printf("New user : %v\n", *user)
	DB.Create(proj)
	// return GetProjByName(&Project.Name), nil
	return nil, nil
}
