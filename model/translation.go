package model

type Translation struct {
	ID      int64  `gorm:"primaryKey;autoIncrement"`
	UserID  int64  `json:"user_id" gorm:"index"`
	ProjID  int64  `json:"project_id" gorm:"not null"`
	FileID  int64  `gorm:"index"`
	SentID  int64  `json:"phrase_id" gorm:"not null"`
	TrnLang string ``
}
