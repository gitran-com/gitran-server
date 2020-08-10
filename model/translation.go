package model

type Translation struct {
	UID        uint64 `gorm:"not_null"`
	PID        uint64 `gorm:"not_null"`
	TargetLang string `gorm:"type:char(3)"`
}
