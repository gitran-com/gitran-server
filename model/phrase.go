package model

type Phrase struct {
	PhrzID uint64 `json:"phrase_id" gorm:"primary_key;auto_increment"`
	ProjID uint64 `json:"project_id" gorm:"not null;index"`
	FileID uint64 `gorm:"not null;index"`
	Text   string `gorm:"type:text"`
}
