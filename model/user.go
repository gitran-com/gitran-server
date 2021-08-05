package model

import (
	"bytes"
	"crypto/sha512"
	"strconv"
	"time"

	"github.com/markbates/goth"
	"github.com/gitran-com/gitran-server/config"
	"github.com/gitran-com/gitran-server/util"
)

//User means user
type User struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"type:varchar(32);index"`
	Email       string    `json:"email" gorm:"type:varchar(256);uniqueIndex"`
	AvatarURL   string    `json:"avatar_url" gorm:"type:varchar(256)"`
	Bio         string    `json:"bio" gorm:"type:varchar(256)"`
	GithubID    int64     `json:"github_id" gorm:"index"`
	IsActive    bool      `json:"is_active" gorm:"index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	LastLoginAt time.Time `json:"-"`
	Password    []byte    `json:"-" gorm:"type:binary(64)"`
	Salt        []byte    `json:"-" gorm:"type:binary(64)"`
}

//TableName return table name
func (*User) TableName() string {
	return config.DB.TablePrefix + "users"
}

//GetUserByID gets a user by id
func GetUserByID(id int64) *User {
	var user []User
	db.First(&user, id)
	if len(user) > 0 {
		return &user[0]
	}
	return nil
}

//GetUserByEmail gets a user by email
func GetUserByEmail(email string) *User {
	var user []User
	db.Where("email=?", email).First(&user)
	if len(user) > 0 {
		return &user[0]
	}
	return nil
}

//GetUserByGithubID gets a user by github id
func GetUserByGithubID(ghid int64) *User {
	var user []User
	db.Where("github_id=?", ghid).First(&user)
	if len(user) > 0 {
		return &user[0]
	}
	return nil
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

//CreateUser creates a new user
func CreateUser(user *User) (*User, error) {
	if res := db.Create(user); res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

//NewUserFromGithub creates a new user from OAuth
func NewUserFromGithub(ext *goth.User) (*User, error) {
	ext_id, _ := strconv.ParseInt(ext.UserID, 10, 64)
	bio, ok := ext.RawData["bio"].(string)
	if !ok {
		bio = ""
	}
	user := &User{
		Name:      ext.Name,
		Email:     ext.Email,
		AvatarURL: ext.AvatarURL,
		Bio:       bio,
		GithubID:  ext_id,
		IsActive:  true,
	}
	return user, nil
}

//GenSalt gens a random 64-byte salt
func GenSalt() string {
	return util.RandString(64)
}

//UpdateUserGithubID update a user github_id
func UpdateUserGithubID(user *User, github_id int64) *User {
	db.Model(user).Select("github_id").Updates(map[string]interface{}{"github_id": github_id})
	user.GithubID = github_id
	return user
}
