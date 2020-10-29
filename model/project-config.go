package model

import "github.com/wzru/gitran-server/config"

//ProjCfg means project config
type ProjCfg struct {
	ID        uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
	FileName  string `json:"filename" gorm:"notNull"`
	ProjID    uint64 `json:"project_id" gorm:"index;notNull"`
	SrcBr     string `json:"src_branch" gorm:"notNull"`
	TgtBr     string `json:"tgt_branch" gorm:"notNull"`
	SyncTime  uint64 `json:"sync_time,omitempty" gorm:"notNull"`
	PushTrans bool   `json:"push_trans" gorm:"notNull"`
}

//TableName return table name
func (*ProjCfg) TableName() string {
	return config.DB.TablePrefix + "project_configs"
}
