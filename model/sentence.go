package model

import (
	"errors"

	"github.com/gitran-com/gitran-server/config"
	"gorm.io/gorm"
)

type Sentence struct {
	ID      int64  `gorm:"primaryKey;autoIncrement"`
	ProjID  int64  `gorm:"index"`
	FileID  int64  `gorm:"index"`
	Valid   bool   `gorm:"index"`
	Content string `gorm:"index"`
}

//TableName return table name
func (*Sentence) TableName() string {
	return config.DB.TablePrefix + "sentences"
}

func (sent *Sentence) Write() {
	db.Save(sent)
}

//NewSentence creates a new file for project
func NewSentence(sent *Sentence) (*Sentence, error) {
	if result := db.Create(sent); result.Error != nil {
		return nil, result.Error
	}
	return sent, nil
}

func MustGetValidSent(proj_id int64, file_id int64, content string) *Sentence {
	sent := &Sentence{}
	res := db.Where("file_id=? AND content=?", file_id, content).First(sent)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		sent, _ = NewSentence(&Sentence{
			ProjID:  proj_id,
			FileID:  file_id,
			Valid:   true,
			Content: content,
		})
		return sent
	}
	sent.Valid = true
	sent.Write()
	return sent
}

func SetAllSentsInvalid(file_id int64) {
	db.Model(&Sentence{}).Where("file_id=?", file_id).Updates(map[string]interface{}{"valid": false})
}
