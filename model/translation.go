package model

type Translation struct {
	UserID     uint64 `json:"user_id" gorm:"not null"`
	ProjID     uint64 `json:"project_id" gorm:"not null"`
	PhrzID     uint64 `json:"phrase_id" gorm:"not null"`
	TargetLang uint16 ``
}
