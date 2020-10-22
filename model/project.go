package model

import (
	"time"

	"github.com/wzru/gitran-server/config"
)

//ProjInfo means project's infomation
type ProjInfo struct {
	ID        uint64        `json:"id" gorm:"primaryKey;autoIncrement"`
	OwnerID   uint64        `json:"owner_id" gorm:"index;notNull"`
	IsUsers   bool          `json:"is_users" gorm:"index;notNull"`
	Name      string        `gorm:"type:varchar(32);index;notNull"`
	Desc      string        `json:"description" gorm:"type:varchar(256)"`
	IsPrivate bool          `json:"is_private"`
	IsGit     bool          `json:"is_git" gorm:"type:notNull"`
	GitURL    string        `json:"git_url,omitempty" gorm:"type:varchar(256)"`
	GitBranch string        `json:"git_branch,omitempty" gorm:"type:varchar(32)"`
	SyncTime  uint64        `json:"sync_time,omitempty"`
	SrcLangs  []config.Lang `json:"source_languages" gorm:"type:varchar(128)"`
	TgtLangs  []config.Lang `json:"target_languages" gorm:"type:varchar(128)"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

//Project means project model
type Project struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	OwnerID   uint64    `json:"owner_id" gorm:"index;notNull"`
	IsUsers   bool      `json:"is_users" gorm:"index;notNull"`
	Name      string    `gorm:"type:varchar(32);index;notNull"`
	Desc      string    `json:"description" gorm:"type:varchar(256)"`
	IsPrivate bool      `json:"is_private" gorm:"notNull"`
	IsGit     bool      `json:"is_git" gorm:"notNull"`
	GitURL    string    `json:"git_url,omitempty" gorm:"type:varchar(256)"`
	GitBranch string    `json:"git_branch,omitempty" gorm:"type:varchar(32)"`
	SyncTime  uint64    `json:"sync_time,omitempty"`
	SrcLangs  string    `json:"source_languages" gorm:"type:varchar(128)"`
	TgtLangs  string    `json:"target_languages" gorm:"type:varchar(128)"`
	Path      string    `gorm:"type:varchar(256)"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//TableName return table name
func (*Project) TableName() string {
	return config.DB.TablePrefix + "projects"
}

//GetProjByOwnerName get project by owner & name
func GetProjByOwnerName(id uint64, name string) *Project {
	var proj []Project
	db.Where("owner_id=? AND name=?", id, name).First(&proj)
	if len(proj) > 0 {
		return &proj[0]
	}
	return nil
}

//NewProj creates a new project
func NewProj(proj *Project) (*Project, error) {
	if result := db.Create(proj); result.Error != nil {
		return nil, result.Error
	}
	return proj, nil
}
