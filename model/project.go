package model

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/gitran-com/gitran-server/config"
	"github.com/gitran-com/gitran-server/constant"
	"github.com/gitran-com/gitran-server/util"
	log "github.com/sirupsen/logrus"
)

var (
	ProjMutexMap = util.NewMutexMap()
)

//Project means project model
type Project struct {
	ID                 int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	URI                string     `json:"uri" gorm:"type:varchar(32);uniqueIndex;notNull"`
	Name               string     `json:"name" gorm:"type:varchar(32);notNull"`
	OwnerID            int64      `json:"owner_id" gorm:"index;notNull"`
	Token              string     `json:"-" gorm:"type:varchar(128)"`
	Type               int        `json:"type" gorm:"type:tinyint;index;notNull"`
	Status             int        `json:"status" gorm:"type:tinyint;notNull"`
	Desc               string     `json:"desc" gorm:"type:varchar(256)"`
	GitURL             string     `json:"git_url" gorm:"type:varchar(256)"`
	Path               string     `json:"-" gorm:"type:varchar(256)"`
	SrcLangs           string     `json:"-" gorm:"type:varchar(128)"`
	TrnLangs           string     `json:"-" gorm:"type:varchar(128)"`
	SourceLanguages    []Language `json:"src_langs" gorm:"-"`
	TranslateLanguages []Language `json:"trn_langs" gorm:"-"`
	ErrMsg             string     `json:"error_message"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
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
	var (
		gitURL string
		err    error
	)
	for i := 0; i < constant.MaxProjInitRetry; i++ {
		if proj.Type == constant.ProjTypeGitURL || proj.Type == constant.ProjTypeGithub {
			if proj.Token == "" {
				gitURL = proj.GitURL
			} else {
				url, _ := util.ParseGitURL(proj.GitURL)
				// fmt.Printf("url=%+v\n", *url)
				gitURL = fmt.Sprintf("https://gitran:%s@%s/%s", proj.Token, url.Host, url.Path)
			}
			// fmt.Printf("GitURL=%v\n", gitURL)
			cmd := exec.Command("git", "clone", "--no-single-branch", "--depth=1", gitURL, proj.Path)
			// fmt.Printf("cmd=%v\n", cmd.String())
			err := cmd.Run()
			if err == nil {
				break
			} else {
				log.Warnf("git clone %v error : %v", proj.URI, err.Error())
			}
		} else if proj.Type == constant.ProjTypePlain {
			break
		} else {
			//TODO
			err = fmt.Errorf("Project.Init error: type %v has not been implemented", proj.Type)
			break
		}
		time.Sleep(time.Second * 5)
	}
	if err != nil {
		proj.InitFail(err)
	} else {
		proj.InitSucc()
		NewProjCfg(&ProjCfg{ID: proj.ID})
	}
}

func (proj *Project) UpdateStatus(stat int) {
	proj.Status = stat
	proj.Write()
}

func (proj *Project) InitFail(err error) {
	proj.Status = constant.ProjStatFailed
	proj.ErrMsg = err.Error()
	proj.Write()
}

func (proj *Project) InitSucc() {
	proj.Status = constant.ProjStatReady
	proj.ErrMsg = ""
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

//ListUninitProj list all uninitialized projects
func ListUninitProj() []Project {
	var projs []Project
	db.Where("status IN(?,?)", constant.ProjStatCreated, constant.ProjStatFailed).Find(&projs)
	// for i := range projs {
	// 	projs[i].FillLangs()
	// }
	return projs
}
