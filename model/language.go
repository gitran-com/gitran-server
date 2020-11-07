package model

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/wzru/gitran-server/constant"
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

func CheckLangs(langs []Language) bool {
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

//GetLangsFromString extract []Lang from string like "eng|zho"
func GetLangsFromString(s string) []Language {
	if s == "" {
		return nil
	}
	ss := strings.Split(s, constant.Delim)
	langs := make([]Language, len(ss))
	for i, code := range ss {
		lang := GetLangByCode(code)
		if lang == nil {
			log.Warnf("unknown language %v", code)
		} else {
			langs[i] = *GetLangByCode(code)
		}
	}
	return langs
}
