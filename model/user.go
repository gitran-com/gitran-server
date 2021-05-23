package model

import (
	"bytes"
	"crypto/sha512"
	"strconv"
	"time"

	"github.com/markbates/goth"
	"github.com/wzru/gitran-server/config"
	"github.com/wzru/gitran-server/util"
)

//User means user
type User struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"type:varchar(32);uniqueIndex"`
	Email       string    `gorm:"type:varchar(64);index"`
	AvatarURL   string    `gorm:"type:varchar(128)"`
	Bio         string    `gorm:"type:varchar(128)"`
	GithubID    uint64    `gorm:"index"`
	IsActive    bool      `gorm:"index"`
	CreatedAt   time.Time ``
	UpdatedAt   time.Time ``
	LastLoginAt time.Time ``
	Password    []byte    `gorm:"type:binary(64)"`
	Salt        []byte    `gorm:"type:binary(64)"`
}

//UserInfo means user's infomation
type UserInfo struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email"`
	AvatarURL string    `json:"avatar_url"`
	Bio       string    `json:"bio"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//TableName return table name
func (*User) TableName() string {
	return config.DB.TablePrefix + "users"
}

//GetUserByNameEmail gets a user by name or email
func GetUserByNameEmail(login string, email string) *User {
	var user []User
	db.Where("name=? OR email=?", login, email).First(&user)
	if len(user) > 0 {
		return &user[0]
	} else {
		return nil
	}
}

//GetUserByName gets a user by login_name
func GetUserByName(name string) *User {
	var user []User
	db.Where("name=?", name).First(&user)
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
func GetUserByGithubID(ghid uint64) *User {
	var user []User
	db.Where("github_id=?", ghid).First(&user)
	if len(user) > 0 {
		return &user[0]
	}
	return nil
}

//GetUserInfoFromUser gen UserInfo from User
func GetUserInfoFromUser(user *User) *UserInfo {
	return &UserInfo{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
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
	ext_id, _ := strconv.ParseUint(ext.UserID, 10, 64)
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
func UpdateUserGithubID(user *User, github_id uint64) *User {
	db.Model(user).Select("github_id").Updates(map[string]interface{}{"github_id": github_id})
	user.GithubID = github_id
	return user
}
