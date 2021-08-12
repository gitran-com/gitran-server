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
	URI                string     `json:"uri" gorm:"type:varchar(32);uniqueIndex;notNull"`
	Name               string     `json:"name" gorm:"type:varchar(32);notNull"`
	OwnerID            int64      `json:"owner_id" gorm:"index;notNull"`
	Token              string     `json:"-"`
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
	SetUserProjRole(proj.OwnerID, proj.ID, RoleAdmin)
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

func (proj *Project) FillLangs() {
	proj.SourceLanguages, _ = ParseLangsFromStr(proj.SrcLangs)
	proj.TranslateLanguages, _ = ParseLangsFromStr(proj.TrnLangs)
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
		proj[0].FillLangs()
		return &proj[0]
	}
	return nil
}

//ListUserProj list all projects from a user
func ListUserProj(user_id int64) []Project {
	var projs []Project
	db.Where("owner_id=?", user_id).Find(&projs)
	for i := range projs {
		projs[i].FillLangs()
	}
	return projs
}
