package model

import (
	"time"

	"github.com/gitran-com/gitran-server/config"
)

//Token 存储Github等来源的Token
type Token struct {
	ID             int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Valid          bool      `json:"-"`
	Source         int       `json:"-" gorm:"type:varchar(16);index"`
	OwnerID        int64     `json:"-" gorm:"index"`
	OwnerName      string    `json:"owner_name" gorm:"type:varchar(32)"`
	OwnerAvatarURL string    `json:"avatar_url" gorm:"type:varchar(128)"`
	AccessToken    string    `json:"-" gorm:"type:varchar(128)"`
	Scope          string    `json:"-" gorm:"type:varchar(8);index"`
	CreatedAt      time.Time `json:"created_at"`
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

//GetValidTokensByOwnerID get a token by owner id
func GetValidTokensByOwnerID(oid int64, src int) []Token {
	var tk []Token
	db.Where("owner_id=? AND source=? AND valid=?", oid, src, true).First(&tk)
	return tk
}

//GetTokenByID gets a token by id
func GetTokenByID(id int64) *Token {
	var tk []Token
	db.First(&tk, id)
	if len(tk) > 0 {
		return &tk[0]
	}
	return nil
}
