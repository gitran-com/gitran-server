package model

import (
	"time"

	"github.com/wzru/gitran-server/config"
)

//Token 存储Github等来源的Token
type Token struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement"`
	Source      int       `gorm:"index"`
	OwnerID     uint64    `gorm:"index"`
	AccessToken string    `json:"access_token" gorm:"type:varchar(128)"`
	TokenType   string    `json:"token_type" gorm:"type:varchar(8)"`
	Scope       string    `json:"scope" gorm:"type:varchar(8);index"`
	CreatedAt   time.Time `json:"created_at"`
}

//TableName return table name
func (*Token) TableName() string {
	return config.DB.TablePrefix + "tokens"
}

//NewToken create a new token
func NewToken(tk *Token) (*Token, error) {
	if res := db.Create(tk); res.Error != nil {
		return nil, res.Error
	}
	return tk, nil
}

//GetTokenByOwnerID get a token by owner id
func GetTokenByOwnerID(oid uint64, src int) *Token {
	var tk []Token
	db.Where("owner_id=? AND source=?", oid, src).First(&tk)
	if len(tk) > 0 {
		return &tk[0]
	}
	return nil
}
