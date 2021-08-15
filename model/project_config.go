package model

import (
	"bytes"
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gitran-com/gitran-server/config"
	"github.com/gitran-com/gitran-server/constant"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"gorm.io/gorm"
)

//ProjCfg means project config
type ProjCfg struct {
	ID            int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	Status        int        `json:"status" gorm:"type:tinyint;"`
	SrcBr         string     `json:"src_branch" gorm:"type:varchar(32);notNull"`
	TrnBr         string     `json:"trn_branch" gorm:"type:varchar(32);notNull"`
	PullGap       uint16     `json:"pull_gap" gorm:"index;notNull"`
	PushGap       uint16     `json:"push_gap" gorm:"index;notNull"`
	PullStatus    int        `json:"pull_status" gorm:"notNull"`
	PushStatus    int        `json:"push_status" gorm:"notNull"`
	LastPullAt    *time.Time `json:"last_pull_at"`
	LastPushAt    *time.Time `json:"last_push_at"`
	FileMapsBytes []byte     `json:"-" gorm:"column:file_maps"`
	FileMaps      []FileMap  `json:"file_maps" gorm:"-"`
	IgnRegsBytes  []byte     `json:"-" gorm:"column:ignores"`
	IgnRegs       []string   `json:"ignores" gorm:"-"`
}

//FileMap means source files => translate maps
type FileMap struct {
	SrcFileReg string `json:"src_files"`
	TrnFileReg string `json:"trn_files"`
}

//TableName return table name
func (*ProjCfg) TableName() string {
	return config.DB.TablePrefix + "project_configs"
}

//Write writes project config to DB
func (cfg *ProjCfg) Write() {
	db.Save(cfg)
}

func (cfg *ProjCfg) AfterFind(tx *gorm.DB) error {
	json.Unmarshal(cfg.FileMapsBytes, &cfg.FileMaps)
	json.Unmarshal(cfg.IgnRegsBytes, &cfg.IgnRegs)
	return nil
}

//NewProjCfg creates a new cfg for project
func NewProjCfg(cfg *ProjCfg) (*ProjCfg, error) {
	if result := db.Create(cfg); result.Error != nil {
		return nil, result.Error
	}
	return cfg, nil
}

func (cfg *ProjCfg) UpdateProjCfg(req *UpdateProjCfgRequest) error {
	var (
		needReprocess bool
	)
	if req.SrcBr != cfg.SrcBr ||
		!bytes.Equal(req.FileMapsBytes, cfg.FileMapsBytes) ||
		!bytes.Equal(req.IgnRegsBytes, cfg.IgnRegsBytes) {
		needReprocess = true
	}
	if err := db.Model(cfg).Updates(req.Map()).Error; err != nil {
		return err
	}
	if needReprocess {
		go cfg.Process()
	}
	return nil
}

func (cfg *ProjCfg) UpdateStatus(stat int) {
	cfg.Status = stat
	cfg.Write()
}

//GetProjCfgByID get a project config by config id
func GetProjCfgByID(id int64) *ProjCfg {
	var pc []ProjCfg
	db.Where("id=?", id).First(&pc)
	if len(pc) > 0 {
		return &pc[0]
	}
	return nil
}

//NewBrchRule create a new branch rule in DB
func NewBrchRule(rule *FileMap) (*FileMap, error) {
	if res := db.Create(rule); res.Error != nil {
		return nil, res.Error
	}
	return rule, nil
}

//UpdateProjCfgPullStatus update a project cfg pull status
func UpdateProjCfgPullStatus(cfg *ProjCfg, stat int) {
	db.Model(cfg).Select("pull_status").Updates(map[string]interface{}{"pull_status": stat})
}

//UpdateProjCfgPushStatus update a project cfg push status
func UpdateProjCfgPushStatus(cfg *ProjCfg, stat int) {
	db.Model(cfg).Select("push_status").Updates(map[string]interface{}{"push_status": stat})
}

//ListSyncProjCfg list all project cfg that should be sync
func ListSyncProjCfg() []ProjCfg {
	var cfg []ProjCfg
	db.Where("pull_gap!=0 OR push_gap!=0").Find(&cfg)
	return cfg
}

func (cfg *ProjCfg) Process() {
	var (
		proj = GetProjByID(cfg.ID)
		stat = constant.ProjStatReady
	)
	lk := ProjMutexMap.Lock(proj.ID)
	defer lk.Unlock()
	defer cfg.UpdateStatus(stat)
	cfg.UpdateStatus(constant.ProjStatProcessing)
	repo, err := git.PlainOpen(proj.Path)
	if err != nil {
		log.Errorf("ProjCfg.Process error when git.PlainOpen(%s): %+v", proj.Path, err.Error())
		return
	}
	wt, _ := repo.Worktree()
	srcBr := "refs/heads/" + cfg.SrcBr
	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(srcBr),
	})
	if err != nil {
		log.Errorf("ProjCfg.Process() error when git.Checkout(%s): %+v", srcBr, err.Error())
		stat = constant.ProjStatFailed
		return
	}
	//TODO
}

func GetFileMapsSrcFiles(fms []FileMap) []string {
	files := []string{}
	for _, fm := range fms {
		files = append(files, fm.SrcFileReg)
	}
	return files
}
