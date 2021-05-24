package model

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/wzru/gitran-server/config"
	"github.com/wzru/gitran-server/constant"
)

//ProjCfg means project config
type ProjCfg struct {
	ID         uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	ProjID     uint64    `json:"project_id" gorm:"index;notNull"`
	FileName   string    `json:"file_name" gorm:"type:varchar(32);notNull"`
	Changed    bool      `json:"-" gorm:"index;notNull"`
	SrcBr      string    `json:"src_branch" gorm:"type:varchar(32);notNull"`
	TrnBr      string    `json:"trn_branch" gorm:"type:varchar(32);notNull"`
	PullItv    uint16    `json:"pull_interval" gorm:"index;notNull"`
	PushItv    uint16    `json:"push_interval" gorm:"index;notNull"`
	PullStatus int       `json:"pull_status" gorm:"notNull"`
	PushStatus int       `json:"push_status" gorm:"notNull"`
	LastPullAt time.Time `json:"last_pull_at"`
	LastPushAt time.Time `json:"last_push_at"`
}

//BrchRule means branch rules
type BrchRule struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	ProjCfgID uint64 `gorm:"index;notNull"`
	Status    int    `gorm:"index;notNull"`
	SrcFiles  string `yaml:"source_files" json:"src_files" gorm:"type:varchar(128);notNull"`
	TrnFiles  string `yaml:"translation_files" json:"trn_files" gorm:"type:varchar(128);notNull"`
	IgnFiles  string `yaml:"ignore_files" json:"ign_files" gorm:"type:varchar(256);"`
	Extension string `gorm:"type:varchar(256)"`
}

//BrchRuleInfo means branch rules info
type BrchRuleInfo struct {
	ID        uint64                 `yaml:"-" json:"id"`
	ProjCfgID uint64                 `yaml:"-" json:"project_config_id"`
	Status    int                    `yaml:"-" json:"status"`
	SrcFiles  string                 `yaml:"source_files" json:"src_files"`
	TrnFiles  string                 `yaml:"translation_files" json:"trn_files"`
	IgnFiles  []string               `yaml:"ignore_files" json:"ign_files"`
	Extension map[string]interface{} `yaml:"extension" json:"extension"`
}

//TableName return table name
func (*ProjCfg) TableName() string {
	return config.DB.TablePrefix + "project_configs"
}

//TableName return table name
func (*BrchRule) TableName() string {
	return config.DB.TablePrefix + "branch_rules"
}

//NewProjCfg creates a new cfg for project
func NewProjCfg(cfg *ProjCfg) (*ProjCfg, error) {
	if result := db.Create(cfg); result.Error != nil {
		return nil, result.Error
	}
	return cfg, nil
}

//ListProjCfgByProjID list all project config by project id
func ListProjCfgByProjID(pid uint64) []ProjCfg {
	var pc []ProjCfg
	db.Where("proj_id=?", pid).Find(&pc)
	return pc
}

//GetProjCfgByID get a project config by config id
func GetProjCfgByID(cid uint64) *ProjCfg {
	var pc []ProjCfg
	db.Where("id=?", cid).First(&pc)
	if len(pc) > 0 {
		return &pc[0]
	}
	return nil
}

//NewBrchRule create a new branch rule in DB
func NewBrchRule(rule *BrchRule) (*BrchRule, error) {
	if res := db.Create(rule); res.Error != nil {
		return nil, res.Error
	}
	return rule, nil
}

//ListBrchRuleByCfgID list all branch rules by project config id
func ListBrchRuleByCfgID(cid uint64) []BrchRule {
	var br []BrchRule
	db.Where("proj_cfg_id=?", cid).Find(&br)
	return br
}

//GetBrchRuleInfoFromBrchRule get a branch rule info from a branch rule
func GetBrchRuleInfoFromBrchRule(rule *BrchRule) *BrchRuleInfo {
	if rule == nil {
		return nil
	}
	var ext map[string]interface{}
	json.Unmarshal([]byte(rule.Extension), &ext)
	return &BrchRuleInfo{
		ID:        rule.ID,
		ProjCfgID: rule.ProjCfgID,
		Status:    rule.Status,
		SrcFiles:  rule.SrcFiles,
		TrnFiles:  rule.TrnFiles,
		IgnFiles:  strings.Split(rule.IgnFiles, constant.Delim),
		Extension: ext,
	}
}

//GetBrchRuleInfosFromBrchRules get many branch rule infos from branch rules
func GetBrchRuleInfosFromBrchRules(rules []BrchRule) []BrchRuleInfo {
	var ri []BrchRuleInfo
	for _, rule := range rules {
		ri = append(ri, *GetBrchRuleInfoFromBrchRule(&rule))
	}
	return ri
}

//UpdateProjCfgChanged update a project cfg changed status
func UpdateProjCfgChanged(cfg *ProjCfg, changed bool) {
	db.Model(cfg).Select("changed").Updates(map[string]interface{}{"changed": changed})
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
	db.Where("pull_itv!=0 OR push_itv!=0").Find(&cfg)
	return cfg
}
