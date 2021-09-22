package model

import (
	"crypto/md5"
	"errors"
	"fmt"

	"github.com/gitran-com/gitran-server/config"
	"gorm.io/gorm"
)

type Sentence struct {
	ID      int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	ProjID  int64  `json:"proj_id" gorm:"index"`
	FileID  int64  `json:"file_id" gorm:"index"`
	SeqNo   int    `json:"seq_no" gorm:"index"`
	Offset  int    `json:"offset" gorm:""`
	Valid   bool   `json:"-" gorm:"index"`
	MD5     string `json:"-" gorm:"index;type:char(32)"`
	Content string `json:"content" gorm:"type:text"`
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

func MustGetValidSent(proj_id int64, file_id int64, seq_no int, content string) *Sentence {
	sent := &Sentence{}
	hash := fmt.Sprintf("%x", md5.Sum([]byte(content)))
	res := db.Where("file_id=? AND md5=?", file_id, hash).First(sent)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		sent, _ = NewSentence(&Sentence{
			ProjID:  proj_id,
			FileID:  file_id,
			SeqNo:   seq_no,
			Valid:   true,
			Content: content,
			MD5:     hash,
		})
		return sent
	}
	sent.SeqNo = seq_no
	sent.Valid = true
	sent.Write()
	return sent
}

func SetAllSentsInvalid(file_id int64) {
	db.Model(&Sentence{}).Where("file_id=?", file_id).Updates(map[string]interface{}{"valid": false})
}
