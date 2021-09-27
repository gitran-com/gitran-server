package model

import "github.com/gitran-com/gitran-server/config"

type Translation struct {
	ID       int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID   int64  `json:"user_id" gorm:"index"`
	UserName string `json:"user_name" gorm:""`
	ProjID   int64  `json:"proj_id" gorm:""`
	FileID   int64  `json:"file_id" gorm:"index"`
	SentID   int64  `json:"sent_id" gorm:"index"`
	LangCode string `json:"lang_code" gorm:"index"`
	Content  string `json:"content" gorm:"type:text"`
}

func (*Translation) TableName() string {
	return config.DB.TablePrefix + "translations"
}

func (tran *Translation) Write() {
	db.Save(tran)
}

func ListSentTrans(lang_code string, sent_id int64) []Translation {
	var trans []Translation
	res := db.Raw("SELECT translations.id, users.id AS user_id, users.name AS user_name, proj_id, file_id, sent_id, content, lang_code FROM translations LEFT JOIN users ON users.id=translations.user_id WHERE sent_id=? AND lang_code=?", sent_id, lang_code).Scan(&trans)
	if res.Error != nil {
		return nil
	}
	return trans
}

func GetTran(sent_id int64, user *User, lang_code string) *Translation {
	var tran Translation
	res := db.Where("sent_id=? AND user_id=? AND lang_code=?", sent_id, user.ID, lang_code).First(&tran)
	if res.Error != nil {
		return nil
	}
	tran.UserName = user.Name
	return &tran
}
