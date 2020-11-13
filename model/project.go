package model

import (
	"github.com/wzru/gitran-server/config"
	"github.com/wzru/gitran-server/constant"
)

//Project means project model
type Project struct {
	ID        uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"type:varchar(32);index;notNull"`
	OwnerID   uint64 `json:"owner_id" gorm:"index;notNull"`
	OwnerType uint8  `json:"owner_type" gorm:"index;notNull"`
	Type      uint8  `json:"type" gorm:"index;notNull"`
	//Status means whether project is ready
	Status    uint8  `json:"status" gorm:"notNull"`
	Desc      string `json:"desc" gorm:"type:varchar(256)"`
	Private   bool   `json:"private" gorm:"index;notNull"`
	GitURL    string `json:"git_url,omitempty" gorm:"type:varchar(256)"`
	RepoID    uint64 ``
	Path      string `gorm:"type:varchar(256)"`
	SrcLangs  string `json:"src_langs" gorm:"type:varchar(128)"`
	TrnLangs  string `json:"trn_langs" gorm:"type:varchar(128)"`
	CreatedAt int64  `gorm:"autoCreateTime:nano"`
	UpdatedAt int64  `gorm:"autoUpdateTime:nano"`
}

//ProjInfo means project's infomation
type ProjInfo struct {
	ID        uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	OwnerID   uint64     `json:"owner_id" gorm:"index;notNull"`
	OwnerType uint8      `json:"owner_type" gorm:"index;notNull"`
	Type      uint8      `json:"type" gorm:"index;notNull"`
	Status    uint8      `json:"status" gorm:"notNull"`
	Name      string     `json:"name" gorm:"type:varchar(32);index;notNull"`
	Desc      string     `json:"desc" gorm:"type:varchar(256)"`
	Private   bool       `json:"private"`
	GitURL    string     `json:"git_url,omitempty" gorm:"type:varchar(256)"`
	SrcLangs  []Language `json:"src_langs"`
	TrnLangs  []Language `json:"trn_langs"`
	CreatedAt int64      `json:"created_at"`
	UpdatedAt int64      `json:"updated_at"`
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
func GetProjByOwnerIDName(oid uint64, name string, priv bool) *Project {
	var proj []Project
	if priv {
		db.Where("owner_id=? AND name=?", oid, name).First(&proj)
	} else {
		db.Where("owner_id=? AND name=? AND private=?", oid, name, false).First(&proj)
	}
	if len(proj) > 0 {
		return &proj[0]
	}
	return nil
}

//NewProj creates a new project
func NewProj(proj *Project) (*Project, error) {
	if res := db.Create(proj); res.Error != nil {
		return nil, res.Error
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
		Private:   proj.Private,
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

//ListProjByUserID list projects from a user
func ListProjByUserID(uid uint64, priv bool) []Project {
	var proj []Project
	if priv {
		db.Where("owner_id=? AND owner_type=?", uid, constant.OwnerUsr).Find(&proj)
	} else {
		db.Where("owner_id=? AND owner_type=? AND private=?", uid, constant.OwnerUsr, false).Find(&proj)
	}
	return proj
}

//ListProjFromOwnerID list all projects from an owner id
func ListProjFromOwnerID(oid uint64, priv bool) []Project {
	var proj []Project
	if priv {
		db.Where("owner_id=?", oid).Find(&proj)
	} else {
		db.Where("owner_id=? AND private=?", oid, false).Find(&proj)
	}
	return proj
}

//UpdateProjStatus update a project status
func UpdateProjStatus(proj *Project, status uint8) {
	db.Model(proj).Select("status").Updates(map[string]interface{}{"status": status})
}

//ListProjByStatus list projects by status
func ListProjByStatus(status uint8) []Project {
	var proj []Project
	db.Where("status=?", status).Find(&proj)
	return proj
}
