package model

import (
	"bytes"
	"crypto/sha512"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/wzru/gitran-server/config"
	"github.com/wzru/gitran-server/util"
)

type UserInfo struct {
	ID          uint64        `json:"user_id"`
	Name        string        `json:"name"`
	Email       string        `json:"email"`
	PreferLangs []config.Lang `json:"prefer_langs"`
	Private     bool          `json:"private"`
}

type User struct {
	ID          uint64    `gorm:"primary_key;auto_increment"`
	Login       string    `gorm:"type:varchar(32);unique_index;not null"`
	Name        string    `gorm:"type:varchar(32)"`
	Email       string    `gorm:"type:varchar(64);unique_index"`
	AvatarURL   string    `gorm:"type:varchar(128)"`
	Bio         string    `gorm:"type:varchar(128)"`
	Password    []byte    `gorm:"type:binary(64);not null"`
	Salt        []byte    `gorm:"type:binary(64);not null"`
	GithubID    uint      `gorm:""`
	PreferLangs string    `gorm:"type:varchar(128)"`
	CreatedAt   time.Time ``
	UpdatedAt   time.Time ``
}

func GetUserByNameEmail(name *string, email *string) *User {
	user := &User{}
	DB.Where("name=? OR email=?", *name, *email).First(user)
	return user
}

func GetUserByName(name *string) *User {
	user := &User{}
	DB.Where("name=?", *name).First(user)
	return user
}

func GetLangByCode(code string) *config.Lang {
	for _, lang := range config.Langs {
		if lang.Code == code {
			return &lang
		}
	}
	return nil
}

func GenPreferLangsFromString(s *string) []config.Lang {
	if s == nil || *s == "" {
		return nil
	}
	ss := strings.Split(*s, ",")
	langs := make([]config.Lang, len(ss))
	for i, code := range ss {
		lang := GetLangByCode(code)
		if lang == nil {
			log.Warnf("Unknown language %v", code)
		} else {
			langs[i] = *GetLangByCode(code)
		}
	}
	return langs
}

func GenUserInfoFromUser(user *User, private bool) *UserInfo {
	return &UserInfo{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		PreferLangs: GenPreferLangsFromString(&user.PreferLangs),
		Private:     private,
	}
}

func HashSalt(pass *string, salt []byte) []byte {
	sum := sha512.Sum512(append([]byte(*pass), salt...))
	return sum[:]
}

func CheckPassword(pass *string, user *User) bool {
	return bytes.Equal(HashSalt(pass, user.Salt), user.Password)
}

func NewUser(user *User) (*User, error) {
	//fmt.Printf("New user : %v\n", *user)
	DB.Create(user)
	return GetUserByNameEmail(&user.Name, &user.Email), nil
}

func GenSalt() string {
	return util.RandStringBytesMaskImprSrcUnsafe(64)
}
