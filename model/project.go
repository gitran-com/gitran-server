package model

import (
	"os"
	"time"

	"github.com/gitran-com/gitran-server/config"
	"github.com/gitran-com/gitran-server/constant"
	"github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	log "github.com/sirupsen/logrus"
)

//Project means project model
type Project struct {
	ID                 int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name               string     `json:"name" gorm:"type:varchar(32);uniqueIndex;notNull"`
	OwnerID            int64      `json:"owner_id" gorm:"index;notNull"`
	Token              string     `json:"-" gorm:"-"`
	Type               int        `json:"type" gorm:"index;notNull"`
	Status             int        `json:"status" gorm:"notNull"`
	Desc               string     `json:"desc" gorm:"type:varchar(256)"`
	GitURL             string     `json:"git_url" gorm:"type:varchar(256)"`
	Path               string     `json:"-" gorm:"type:varchar(256)"`
	SrcLangs           string     `json:"-" gorm:"type:varchar(128)"`
	TrnLangs           string     `json:"-" gorm:"type:varchar(128)"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	SourceLanguages    []Language `json:"src_langs" gorm:"-"`
	TranslateLanguages []Language `json:"trn_langs" gorm:"-"`
}

//TableName return table name
func (*Project) TableName() string {
	return config.DB.TablePrefix + "projects"
}

//Write writes project to DB
func (proj *Project) Write() {
	db.Save(proj)
}

//Create creates new project
func (proj *Project) Create() error {
	return db.Create(proj).Error
}

//Init init a new project
func (proj *Project) Init() {
	if proj.Type == constant.ProjTypeGitURL {
		_, err := git.PlainClone(proj.Path, false, &git.CloneOptions{
			URL:          proj.GitURL,
			Progress:     os.Stdout, //TODO: update progress in DB
			Depth:        1,
			SingleBranch: false,
		})
		if err == nil {
			proj.Status = constant.ProjStatReady
			proj.Write()
		} else {
			log.Warnf("git clone error : %v", err.Error())
		}
	} else if proj.Type == constant.ProjTypeGithub {
		_, err := git.PlainClone(proj.Path, false, &git.CloneOptions{
			URL: proj.GitURL,
			Auth: &githttp.BasicAuth{
				Username: config.APP.Name,
				Password: proj.Token,
			},
			Progress:     os.Stdout,
			Depth:        1,
			SingleBranch: false,
		})
		if err == nil {
			proj.Status = constant.ProjStatReady
			proj.Write()
		} else {
			log.Warnf("git clone error : %v", err.Error())
		}
	} else {
		//TODO
		log.Errorf("Project.Init error: type %v has not been implemented", proj.Type)
		return
	}
}

func (proj *Project) UpdateStatus(stat int) {
	proj.Status = stat
	proj.Write()
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
