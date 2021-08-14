package model

import (
	"encoding/json"
	"time"

	"github.com/gitran-com/gitran-server/config"
	"gorm.io/gorm"
)

//ProjCfg means project config
type ProjCfg struct {
	ID               int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	SrcBr            string    `json:"src_branch" gorm:"type:varchar(32);notNull"`
	TrnBr            string    `json:"trn_branch" gorm:"type:varchar(32);notNull"`
	PullGap          uint16    `json:"pull_gap" gorm:"index;notNull"`
	PushGap          uint16    `json:"push_gap" gorm:"index;notNull"`
	PullStatus       int       `json:"pull_status" gorm:"notNull"`
	PushStatus       int       `json:"push_status" gorm:"notNull"`
	LastPullAt       time.Time `json:"last_pull_at"`
	LastPushAt       time.Time `json:"last_push_at"`
	PublicView       bool      `json:"public_view"`
	PublicContribute bool      `json:"public_contribute"`
	FileMapsBytes    []byte    `json:"-" gorm:"column:file_maps"`
	FileMaps         []FileMap `gorm:"-"`
	IgnRegsBytes     []byte    `json:"-" gorm:"column:ignores"`
	IgnRegs          []string  `gorm:"-"`
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

func (cfg *ProjCfg) UpdateProjCfg(mp map[string]interface{}) error {
	return db.Model(cfg).Updates(mp).Error
}

//GetCfgByProjID list all project config by project id
func GetCfgByProjID(pid int64) *ProjCfg {
	var pc []ProjCfg
	db.Where("id=?", pid).Find(&pc)
	return &pc[0]
}

//GetProjCfgByID get a project config by config id
func GetProjCfgByID(cid int64) *ProjCfg {
	var pc []ProjCfg
	db.Where("id=?", cid).First(&pc)
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
