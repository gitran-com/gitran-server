package model

import (
	"bytes"
	"crypto/sha512"
	"time"

	"github.com/wzru/gitran-server/config"
	"github.com/wzru/gitran-server/util"
)

//UserInfo means user's infomation
type UserInfo struct {
	ID          uint64        `json:"id"`
	Login       string        `json:"login"`
	Name        string        `json:"name,omitempty"`
	Email       string        `json:"email"`
	AvatarURL   string        `json:"avatar_url"`
	Bio         string        `json:"bio"`
	PreferLangs []config.Lang `json:"prefer_langs"`
	GithubID    uint64        `json:"github_id,omitempty"`
	IsPrivate   bool          `json:"is_private"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

//User means user model
type User struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement"`
	Login       string    `gorm:"type:varchar(32);uniqueIndex;notNull"`
	Name        string    `gorm:"type:varchar(32);index"`
	Email       string    `gorm:"type:varchar(64);uniqueIndex;notNull"`
	AvatarURL   string    `gorm:"type:varchar(128)"`
	Bio         string    `gorm:"type:varchar(128)"`
	GithubID    uint64    `gorm:""`
	PreferLangs string    `gorm:"type:varchar(128)"`
	CreatedAt   time.Time ``
	UpdatedAt   time.Time ``
	Password    []byte    `gorm:"type:binary(64);notNull"`
	Salt        []byte    `gorm:"type:binary(64);notNull"`
}

//TableName return table name
func (*User) TableName() string {
	return config.DB.TablePrefix + "users"
}

//GetUserByLoginEmail gets a user by login or email
func GetUserByLoginEmail(login string, email string) *User {
	var user []User
	db.Where("login=? OR email=?", login, email).First(&user)
	if len(user) > 0 {
		return &user[0]
	} else {
		return nil
	}
}

//GetUserByLogin gets a user by login
func GetUserByLogin(name string) *User {
	var user []User
	db.Where("login=?", name).First(&user)
	if len(user) > 0 {
		return &user[0]
	}
	return nil
}

//GetUserByID gets a user by id
func GetUserByID(id uint64) *User {
	var user []User
	db.First(&user, id)
	if len(user) > 0 {
		return &user[0]
	} else {
		return nil
	}
}

//GetUserInfoFromUser gen UserInfo from User
func GetUserInfoFromUser(user *User, priv bool) *UserInfo {
	if priv {
		return &UserInfo{
			ID:          user.ID,
			Login:       user.Login,
			Name:        user.Name,
			Email:       user.Email,
			PreferLangs: GetLangsFromString(user.PreferLangs),
			GithubID:    user.GithubID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			IsPrivate:   priv,
		}
	} else {
		return &UserInfo{
			ID:          user.ID,
			Login:       user.Login,
			Name:        user.Name,
			Email:       user.Email,
			PreferLangs: GetLangsFromString(user.PreferLangs),
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			IsPrivate:   priv,
		}
	}
}

//HashSalt calcs H(pass+salt)
func HashSalt(pass string, salt []byte) []byte {
	sum := sha512.Sum512(append([]byte(pass), salt...))
	return sum[:]
}

//CheckPasswordCorrect checks whether a password is correct
func CheckPasswordCorrect(pass string, user *User) bool {
	// fmt.Printf("pass=%v\n", pass)
	return bytes.Equal(HashSalt(pass, user.Salt), user.Password)
}

//NewUser creates a new user
func NewUser(user *User) (*User, error) {
	if result := db.Create(user); result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

//GenSalt gens a random 64-byte salt
func GenSalt() string {
	return util.RandStringBytesMaskImprSrcUnsafe(64)
}
