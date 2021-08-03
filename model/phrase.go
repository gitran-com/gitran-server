package model

type Phrase struct {
	PhrzID int64  `json:"phrase_id" gorm:"primary_key;auto_increment"`
	ProjID int64  `json:"project_id" gorm:"not null;index"`
	FileID int64  `gorm:"not null;index"`
	Text   string `gorm:"type:text"`
}
