package model

import (
	"github.com/gitran-com/gitran-server/config"
	"gorm.io/gorm"
)

type Translation struct {
	ID       int64   `gorm:"primaryKey;autoIncrement"`
	UserID   int64   `gorm:"index"`
	ProjID   int64   `gorm:""`
	FileID   int64   `gorm:"index"`
	SentID   int64   `gorm:"index"`
	Likes    int64   `gorm:""`
	Unlikes  int64   `json:"-" gorm:""`
	Score    float64 `json:"-" gorm:"index"`
	LangCode string  `gorm:"index;type:varchar(8)"`
	Pinned   bool    `json:"pinned" gorm:"index"`
	Content  string  `gorm:"type:text"`
}

type TranRes struct {
	ID       int64  `json:"id"`
	UserID   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	Vote     int    `json:"vote"`
	ProjID   int64  `json:"proj_id"`
	FileID   int64  `json:"file_id"`
	SentID   int64  `json:"sent_id"`
	LangCode string `json:"lang_code"`
	Likes    int64  `json:"likes"`
	Content  string `json:"content"`
}

func (*Translation) TableName() string {
	return config.DB.TablePrefix + "translations"
}

func (tran *Translation) Write() {
	db.Save(tran)
}

func GetTranByID(id int64) *Translation {
	var tran Translation
	res := db.First(&tran, id)
	if res.Error != nil {
		return nil
	}
	return &tran
}

func ListSentTrans(lang_code string, sent_id int64) []TranRes {
	var trans []TranRes
	res := db.Raw("SELECT translations.id, users.id AS user_id, users.name AS user_name, vote, proj_id, file_id, sent_id, content, lang_code, likes FROM translations LEFT JOIN users ON users.id=translations.user_id LEFT JOIN votings ON users.id=votings.user_id AND translations.id=votings.tran_id WHERE sent_id=? AND lang_code=? ORDER BY score DESC, unlikes ASC", sent_id, lang_code).Scan(&trans)
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

func (tran *Translation) Delete() {
	db.Transaction(func(tx *gorm.DB) error {
		tx.Where("tran_id=?", tran.ID).Delete(&Voting{})
		tx.Delete(tran)
		return nil
	})
}

func (tran *Translation) Pin() {
	tran.Pinned = true
	tran.Write()
}

func (tran *Translation) Unpin() {
	tran.Pinned = false
	tran.Write()
}
