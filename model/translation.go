package model

type Translation struct {
	UserID     int64  `json:"user_id" gorm:"not null"`
	ProjID     int64  `json:"project_id" gorm:"not null"`
	PhrzID     int64  `json:"phrase_id" gorm:"not null"`
	TargetLang uint16 ``
}
