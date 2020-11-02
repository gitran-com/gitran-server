package model

import "github.com/wzru/gitran-server/config"

//ProjCfg means project config
type ProjCfg struct {
	ID         uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
	ProjID     uint64 `json:"project_id" gorm:"index;notNull"`
	SrcBr      string `json:"src_branch" gorm:"type:varchar(32);notNull"`
	TgtBr      string `json:"tgt_branch" gorm:"type:varchar(32);notNull"`
	SyncTime   uint64 `json:"sync_time,omitempty" gorm:"notNull"`
	SyncStatus bool   `json:"sync_status" gorm:"notNull"`
	PushTrans  bool   `json:"push_trans" gorm:"notNull"`
	FileName   string `json:"filename" gorm:"type:varchar(32);notNull"`
}

//BrchRule means branch rules
type BrchRule struct {
	ID        uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
	ProjCfgID uint64 `json:"project_config_id" gorm:"index;notNull"`
	SrcFiles  string `json:"src_files" gorm:"type:varchar(32);notNull"`
	TrnFiles  string `json:"trn_files" gorm:"type:varchar(32);notNull"`
	Extension string `json:"extension" gorm:"type:varchar(128)"`
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
func NewProjCfg(pc *ProjCfg) (*ProjCfg, error) {
	if result := db.Create(pc); result.Error != nil {
		return nil, result.Error
	}
	return pc, nil
}
