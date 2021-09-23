package model

import "github.com/gitran-com/gitran-server/config"

type Translation struct {
	ID       int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID   int64  `json:"user_id" gorm:"index"`
	ProjID   int64  `json:"proj_id" gorm:""`
	FileID   int64  `json:"file_id" gorm:"index"`
	SentID   int64  `json:"sent_id" gorm:"index"`
	Content  string `json:"content" gorm:"type:text"`
	LangCode string `json:"lang_code" gorm:"index"`
}

func (*Translation) TableName() string {
	return config.DB.TablePrefix + "translations"
}

func (tran *Translation) Write() {
	db.Save(tran)
}

func ListSentTrans(sent_id int64) []Translation {
	var trans []Translation
	res := db.Where("sent_id=?", sent_id).First(&trans)
	if res.Error != nil {
		return nil
	}
	return trans
}

func GetTran(sent_id int64, user_id int64, lang_code string) *Translation {
	var tran Translation
	res := db.Where("sent_id=? AND user_id=? AND lang_code=?", sent_id, user_id, lang_code).First(&tran)
	if res.Error != nil {
		return nil
	}
	return &tran
}
