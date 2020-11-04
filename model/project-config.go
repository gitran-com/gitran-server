package model

import (
	"encoding/json"
	"strings"

	"github.com/wzru/gitran-server/config"
	"github.com/wzru/gitran-server/constant"
)

//ProjCfg means project config
type ProjCfg struct {
	ID         uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
	ProjID     uint64 `json:"project_id" gorm:"index;notNull"`
	Changed    bool   `gorm:"index;notNull"`
	SrcBr      string `json:"src_branch" gorm:"type:varchar(32);notNull"`
	TrnBr      string `json:"trn_branch" gorm:"type:varchar(32);notNull"`
	SyncTime   uint64 `json:"sync_time,omitempty" gorm:"notNull"`
	SyncStatus bool   `json:"sync_status" gorm:"notNull"`
	PushTrans  bool   `json:"push_trans" gorm:"notNull"`
	FileName   string `json:"file_name" gorm:"type:varchar(32);notNull"`
}

//BrchRule means branch rules
type BrchRule struct {
	ID        uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
	ProjCfgID uint64 `json:"project_config_id" gorm:"index;notNull"`
	Status    uint8  `gorm:"index;notNull"`
	SrcFiles  string `yaml:"source_files" json:"src_files" gorm:"type:varchar(128);notNull"`
	TrnFiles  string `yaml:"translation_files" json:"trn_files" gorm:"type:varchar(128);notNull"`
	IgnFiles  string `yaml:"ignore_files" json:"ign_files" gorm:"type:varchar(256);"`
	Extension string `json:"extension" gorm:"type:varchar(256)"`
}

//BrchRuleInfo means branch rules info
type BrchRuleInfo struct {
	ID        uint64                 `yaml:"-" json:"id"`
	ProjCfgID uint64                 `yaml:"-" json:"project_config_id"`
	Status    uint8                  `yaml:"-" json:"status"`
	SrcFiles  string                 `yaml:"source_files" json:"src_files"`
	TrnFiles  string                 `yaml:"translation_files" json:"trn_files"`
	IgnFiles  []string               `yaml:"ignore_files" json:"ign_files"`
	Extension map[string]interface{} `yaml:"extension" json:"extension"`
}

// //NeedUpdate return whether the config need update config file
// func (*ProjCfg) NeedUpdate() bool {
// 	return config.DB.TablePrefix + "project_configs"
// }

//TableName return table name
func (*ProjCfg) TableName() string {
	return config.DB.TablePrefix + "project_configs"
}

//TableName return table name
func (*BrchRule) TableName() string {
	return config.DB.TablePrefix + "branch_rules"
}

//NewProjCfg creates a new cfg for project
func NewProjCfg(pc *ProjCfg) (*ProjCfg, error) {
	if result := db.Create(pc); result.Error != nil {
		return nil, result.Error
	}
	return pc, nil
}

//ListProjCfgFromProjID list all project config by project id
func ListProjCfgFromProjID(pid uint64) []ProjCfg {
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
