package model

import (
	"crypto/sha512"

	"github.com/wzru/gitran-server/config"
	"github.com/wzru/gitran-server/constant"
	"github.com/wzru/gitran-server/util"
)

//User means user
type User struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement"`
	Login       string `gorm:"type:varchar(32);uniqueIndex;notNull"`
	Name        string `gorm:"type:varchar(32);index;notNull"`
	Email       string `gorm:"type:varchar(64);uniqueIndex;notNull"`
	AvatarURL   string `gorm:"type:varchar(128)"`
	Bio         string `gorm:"type:varchar(128)"`
	GithubID    uint64 `gorm:"index"`
	PreferLangs string `gorm:"type:varchar(128)"`
	Salt        string `gorm:"type:bytes;size:64;notNull"`
	Password    string `gorm:"type:bytes;size:64;notNull"`
	CreatedAt   int64  `gorm:"autoCreateTime:nano"`
	UpdatedAt   int64  `gorm:"autoUpdateTime:nano"`
}

//UserInfo means user's infomation
type UserInfo struct {
	ID          uint64     `json:"id"`
	Login       string     `json:"login"`
	Name        string     `json:"name,omitempty"`
	Email       string     `json:"email"`
	AvatarURL   string     `json:"avatar_url"`
	Bio         string     `json:"bio"`
	PreferLangs []Language `json:"prefer_langs"`
	GithubID    uint64     `json:"github_id,omitempty"`
	Private     bool       `json:"private"`
	CreatedAt   int64      `json:"created_at"`
	UpdatedAt   int64      `json:"updated_at"`
}

//TableName return table name
func (*User) TableName() string {
	return config.DB.TablePrefix + "users"
}

//GetUserByNameEmail gets a user by login or email
func GetUserByNameEmail(login string, email string) *User {
	var user []User
	db.Where("login=? OR email=?", login, email).First(&user)
	if len(user) > 0 {
		return &user[0]
	} else {
		return nil
	}
}

//GetUserByName gets a user by login(user)
func GetUserByName(name string) *User {
	var user []User
	db.Where("login=?", name).First(&user)
	if len(user) > 0 {
		return &user[0]
	}
	return nil
}

//GetOwnerByName gets a user or an org by name
func GetOwnerByName(name string) (*User, *Organization, uint8) {
	//TODO
	user := GetUserByName(name)
	if user == nil {
		return nil, nil, constant.OwnerNone
	}
	return user, nil, constant.OwnerUsr
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
			Private:     priv,
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
			Private:     priv,
		}
	}
}

//HashWithSalt calcs H(pass+salt)
func HashWithSalt(pass string, salt string) string {
	sum := sha512.Sum512([]byte(pass + salt))
	return string(sum[:])
}

//CheckPasswordCorrect checks whether a password is correct
func CheckPasswordCorrect(pass string, user *User) bool {
	return HashWithSalt(pass, user.Salt) == user.Password
}

//NewUser creates a new user
func NewUser(user *User) (*User, error) {
	if res := db.Create(user); res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

//GenSalt gens a random 64-byte salt
func GenSalt() string {
	return util.RandString(64)
}

//UpdateUserGithubID update a user github_id
func UpdateUserGithubID(user *User, ghid uint64) {
	db.Model(user).Select("github_id").Updates(map[string]interface{}{"github_id": ghid})
}
