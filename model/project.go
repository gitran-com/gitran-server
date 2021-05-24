package model

import (
	"time"

	"github.com/wzru/gitran-server/config"
	"github.com/wzru/gitran-server/constant"
)

//Project means project model
type Project struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"type:varchar(32);index;notNull"`
	OwnerID   uint64    `gorm:"index;notNull"`
	TokenID   uint64    `gorm:"index"`
	Type      int       `gorm:"index;notNull"`
	Status    int       `gorm:"notNull"`
	Desc      string    `gorm:"type:varchar(256)"`
	GitURL    string    `gorm:"type:varchar(256)"`
	Path      string    `gorm:"type:varchar(256)"`
	SrcLangs  string    `gorm:"type:varchar(128)"`
	TrnLangs  string    `gorm:"type:varchar(128)"`
	CreatedAt time.Time ``
	UpdatedAt time.Time ``
}

//ProjInfo means project's infomation
type ProjInfo struct {
	ID        uint64     `json:"id"`
	Name      string     `json:"name"`
	OwnerID   uint64     `json:"owner_id"`
	Type      int        `json:"type"`
	Status    int        `json:"status"`
	Desc      string     `json:"desc"`
	GitURL    string     `json:"git_url,omitempty"`
	SrcLangs  []Language `json:"src_langs"`
	TrnLangs  []Language `json:"trn_langs"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

//TableName return table name
func (*Project) TableName() string {
	return config.DB.TablePrefix + "projects"
}

//GetProjByID get project by id
func GetProjByID(id uint64) *Project {
	var proj []Project
	db.Where("id=?", id).First(&proj)
	if len(proj) > 0 {
		return &proj[0]
	}
	return nil
}

//GetProjByOwnerIDName get project by oid & name
func GetProjByOwnerIDName(oid uint64, name string, self bool) *Project {
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
		Name:      proj.Name,
		Desc:      proj.Desc,
		Type:      proj.Type,
		Status:    proj.Status,
		GitURL:    proj.GitURL,
		SrcLangs:  GetLangsFromString(proj.SrcLangs),
		TrnLangs:  GetLangsFromString(proj.TrnLangs),
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
	var proj []Project
	if priv {
		db.Where("owner_id=? AND owner_type=?", user.ID, constant.OwnerUsr).Find(&proj)
	} else {
		db.Where("owner_id=? AND owner_type=? AND is_private=?", user.ID, constant.OwnerUsr, false).Find(&proj)
	}
	return proj
}

//ListProjFromOwnerID list all projects from an owner id
func ListProjFromOwnerID(oid uint64, priv bool) []Project {
	var proj []Project
	if priv {
		db.Where("owner_id=?", oid).Find(&proj)
	} else {
		db.Where("owner_id=? AND is_private=?", oid, false).Find(&proj)
	}
	return proj
}

//UpdateProjStatus update a project status
func UpdateProjStatus(proj *Project, status int) {
	db.Model(proj).Select("status").Updates(map[string]interface{}{"status": status})
}
