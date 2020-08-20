package model

import (
	"bytes"
	"crypto/sha512"
	"time"

	"github.com/WangZhengru/gitran-be/util"
)

type User struct {
	ID          uint64    `gorm:"primary_key;auto_increment"`
	Name        string    `gorm:"unique_index;not_null;size:255"`
	Email       string    `gorm:"unique_index;size:255"`
	Password    []byte    `gorm:"type:binary(64);not_null;"`
	Salt        []byte    `gorm:"type:binary(64);not_null;"`
	GithubID    string    ``
	PreferLangs string    `gorm:"type:varchar(128)"`
	CreatedAt   time.Time ``
	UpdatedAt   time.Time ``
}

func GetUserByNameEmail(name *string, email *string) *User {
	user := &User{}
	DB.Where("name=? OR email=?", *name, *email).First(user)
	return user
}

func HashSalt(pass *string, salt []byte) []byte {
	sum := sha512.Sum512(append([]byte(*pass), salt...))
	return sum[:]
}

func CheckPassword(pass *string, user *User) bool {
	return bytes.Equal(HashSalt(pass, user.Salt), user.Password)
}

func NewUser(user *User) (*User, error) {
	DB.Create(user)
	return GetUserByNameEmail(&user.Name, &user.Email), nil
}

func GenSalt() string {
	return util.RandStringBytesMaskImprSrcUnsafe(128)
}
