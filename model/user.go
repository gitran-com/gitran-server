package model

import (
	"bytes"
	"crypto/sha512"
	"errors"
	"strconv"
	"time"

	"github.com/gitran-com/gitran-server/config"
	"github.com/gitran-com/gitran-server/util"
	"github.com/markbates/goth"
	"gorm.io/gorm"
)

//User means user
type User struct {
	ID              int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Email           string    `json:"email" gorm:"type:varchar(256);uniqueIndex"`
	Name            string    `json:"name" gorm:"type:varchar(32);index;notNull"`
	Bio             string    `json:"bio" gorm:"type:varchar(256)"`
	AvatarURL       string    `json:"avatar_url" gorm:"type:varchar(256)"`
	GithubID        int64     `json:"github_id" gorm:"index"`
	GithubRepoToken string    `json:"-" gorm:"type:varchar(64)"`
	IsActive        bool      `json:"is_active"`
	NoPassword      bool      `json:"-"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Password        []byte    `json:"-" gorm:"type:binary(64)"`
	Salt            []byte    `json:"-" gorm:"type:binary(64)"`
}

//TableName returns table name
func (*User) TableName() string {
	return config.DB.TablePrefix + "users"
}

//Write writes user to DB
func (user *User) Write() {
	db.Save(user)
}

//Create creates new user
func (user *User) Create() error {
	return db.Create(user).Error
}

func (user *User) UpdateProfile(mp map[string]interface{}) error {
	return db.Model(user).Updates(mp).Error
}

//GetUserByID gets a user by id
func GetUserByID(id int64) *User {
	var user User
	res := db.First(&user, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	return &user
}

//GetUserByEmail gets a user by email
func GetUserByEmail(email string) *User {
	var user User
	res := db.Where("email=?", email).First(&user)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	return &user
}

//GetUserByGithubID gets a user by github id
func GetUserByGithubID(ghid int64) *User {
	var user User
	res := db.Where("github_id=?", ghid).First(&user)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	return &user
}

//HashSalt calcs H(pass+salt)
func HashSalt(pass string, salt []byte) []byte {
	sum := sha512.Sum512(append([]byte(pass), salt...))
	return sum[:]
}

//CheckPass checks whether a password is correct
func CheckPass(user *User, pass string) bool {
	if user == nil {
		return false
	}
	return bytes.Equal(HashSalt(pass, user.Salt), user.Password)
}

//NewUserFromGithub creates a new user from OAuth
func NewUserFromGithub(ext *goth.User) *User {
	ext_id, _ := strconv.ParseInt(ext.UserID, 10, 64)
	bio, _ := ext.RawData["bio"].(string)
	user := &User{
		Name:       ext.Name,
		Email:      ext.Email,
		AvatarURL:  ext.AvatarURL,
		Bio:        bio,
		GithubID:   ext_id,
		IsActive:   true,
		NoPassword: true,
	}
	return user
}

//GenSalt gens a random 64-byte salt
func GenSalt() []byte {
	return []byte(util.RandString(64))
}
