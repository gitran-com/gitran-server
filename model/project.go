package model

import (
	"time"

	"github.com/wzru/gitran-server/config"
)

//ProjInfo means project's infomation
type ProjInfo struct {
	ID        uint64        `json:"id" gorm:"primaryKey;autoIncrement"`
	OwnerID   uint64        `json:"owner_id" gorm:"index;notNull"`
	OwnerType uint8         `json:"owner_type" gorm:"index;notNull"`
	Type      uint8         `json:"type" gorm:"index;notNull"`
	Name      string        `json:"name" gorm:"type:varchar(32);index;notNull"`
	Desc      string        `json:"desc" gorm:"type:varchar(256)"`
	IsPrivate bool          `json:"is_private"`
	GitURL    string        `json:"git_url,omitempty" gorm:"type:varchar(256)"`
	SyncTime  uint64        `json:"sync_time,omitempty"`
	SrcLangs  []config.Lang `json:"src_langs"`
	TgtLangs  []config.Lang `json:"tgt_langs"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

//Project means project model
type Project struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"type:varchar(32);index;notNull"`
	OwnerID   uint64    `json:"owner_id" gorm:"index;notNull"`
	OwnerType uint8     `json:"owner_type" gorm:"index;notNull"`
	Type      uint8     `json:"type" gorm:"index;notNull"`
	Desc      string    `json:"desc" gorm:"type:varchar(256)"`
	IsPrivate bool      `json:"is_private" gorm:"notNull"`
	GitURL    string    `json:"git_url,omitempty" gorm:"type:varchar(256)"`
	SyncTime  uint64    `json:"sync_time,omitempty"`
	SrcLangs  string    `json:"src_langs" gorm:"type:varchar(128)"`
	TgtLangs  string    `json:"tgt_langs" gorm:"type:varchar(128)"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//TableName return table name
func (*Project) TableName() string {
	return config.DB.TablePrefix + "projects"
}

//GetProjByOIDName get project by oid & name
func GetProjByOIDName(oid uint64, name string, self bool) *Project {
	var proj []Project
	if self {
		db.Where("owner_id=? AND name=?", oid, name).First(&proj)
	} else {
		db.Where("owner_id=? AND name=? AND is_private=?", oid, name, false).First(&proj)
	}
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

//GetProjInfoFromProj get a project info from a project
func GetProjInfoFromProj(proj *Project) *ProjInfo {
	if proj == nil {
		return nil
	}
	return &ProjInfo{
		ID:        proj.ID,
		OwnerID:   proj.OwnerID,
		OwnerType: proj.OwnerType,
		Name:      proj.Name,
		Desc:      proj.Desc,
		IsPrivate: proj.IsPrivate,
		Type:      proj.Type,
		GitURL:    proj.GitURL,
		SyncTime:  proj.SyncTime,
		SrcLangs:  GetLangsFromString(proj.SrcLangs),
		TgtLangs:  GetLangsFromString(proj.TgtLangs),
		CreatedAt: proj.CreatedAt,
		UpdatedAt: proj.UpdatedAt,
	}
}

//GetProjInfosFromProjs get projects info from []project
func GetProjInfosFromProjs(projs []Project) []ProjInfo {
	var pi []ProjInfo
	for _, proj := range projs {
		pi = append(pi, *GetProjInfoFromProj(&proj))
	}
	return pi
}

//ListProjFromUser list all projects from a user
func ListProjFromUser(user *User, priv bool) []Project {
	if user == nil {
		return nil
	}
	return ListProjFromOID(user.ID, priv)
}

//ListProjFromOID list all projects from an owner id
func ListProjFromOID(oid uint64, priv bool) []Project {
	var proj []Project
	if priv {
		db.Where("owner_id=?", oid).First(&proj)
	} else {
		db.Where("owner_id=? AND is_private=?", oid, false).First(&proj)
	}
	return proj
}
