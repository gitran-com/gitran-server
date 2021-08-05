package model

import (
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/gitran-com/gitran-server/constant"
)

//Language 语言
type Language struct {
	ID    int    `json:"id"`
	Code  string `json:"code"`
	Code3 string `json:"code3"`
	ISO   string `json:"iso"`
	Name  string `json:"name"`
}

func GetLangs() []Language {
	return langs
}

func CheckLang(lang *Language) bool {
	for _, cfgLang := range langs {
		if cfgLang.Code == lang.Code {
			return true
		}
	}
	return false
}

func ValidateLangs(langs []Language) bool {
	for _, lang := range langs {
		ok := CheckLang(&lang)
		if !ok {
			return false
		}
	}
	return true
}

//GetLangByID gets a language by language id
func GetLangByID(id int) *Language {
	if id < len(langs) {
		if langs[id].ID == id {
			return &(langs[id])
		}
	}
	for _, lang := range langs {
		if lang.ID == id {
			return &lang
		}
	}
	return nil
}

//GetLangByCode gets a language by language code
func GetLangByCode(code string) *Language {
	for _, lang := range langs {
		if lang.Code == code {
			return &lang
		}
	}
	return nil
}

//ParseLangs extract []Lang from string like "eng|zho"
func ParseLangs(s string) []Language {
	if s == "" {
		return nil
	}
	ss := strings.Split(s, constant.Delim)
	langs := make([]Language, len(ss))
	for i, id := range ss {
		num, _ := strconv.Atoi(id)
		lang := GetLangByID(num)
		if lang == nil {
			log.Warnf("unknown language #%v", num)
		} else {
			langs[i] = *lang
		}
	}
	return langs
}
