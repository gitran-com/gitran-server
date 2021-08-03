package model

import (
	"time"

	"github.com/wzru/gitran-server/config"
)

//Project means project model
type Project struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"type:varchar(32);uniqueIndex;notNull"`
	OwnerID   int64     `gorm:"index;notNull"`
	TokenID   int64     `gorm:"index"`
	Token     *Token    `gorm:"-"`
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
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	URI       string     `json:"uri"`
	OwnerID   int64      `json:"owner_id"`
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
func GetProjByID(id int64) *Project {
	var proj []Project
	db.Where("id=?", id).First(&proj)
	if len(proj) > 0 {
		return &proj[0]
	}
	return nil
}

//GetProjByURI get project by name
func GetProjByURI(uri string) *Project {
	var proj []Project
	db.Where("uri=?", uri).First(&proj)
	if len(proj) > 0 {
		return &proj[0]
	}
	return nil
}

//GetProjByOwnerIDName get project by oid & name
func GetProjByOwnerIDName(oid int64, name string, self bool) *Project {
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
		SrcLangs:  ParseLangs(proj.SrcLangs),
		TrnLangs:  ParseLangs(proj.TrnLangs),
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

//ListUserProj list all projects from a user
func ListUserProj(user_id int64) []Project {
	var proj []Project
	db.Where("owner_id=?", user_id).Find(&proj)
	return proj
}

//ListProjFromOwnerID list all projects from an owner id
func ListProjFromOwnerID(oid int64, priv bool) []Project {
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
