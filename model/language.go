package model

import (
	"strings"

	"github.com/gitran-com/gitran-server/constant"
)

//Language 语言
type Language struct {
	Code  string `json:"code"`
	Code2 string `json:"code2"`
	ISO   string `json:"iso"`
	Name  string `json:"name"`
}

func ListLangs() []Language {
	return langs
}

//GetLangByCode gets a language by language code
func GetLangByCode(code string) (lang Language, ok bool) {
	lang, ok = langMap[code]
	return
}

//ParseLangsFromStr extract []Lang from string like "eng|zho"
func ParseLangsFromStr(s string) ([]Language, bool) {
	if s == "" {
		return nil, true
	}
	return ParseLangsFromCodes(strings.Split(s, constant.Delim))
}

//ParseLangsFromCodes extract []Lang from []string like ["en", "zh-CN"]
func ParseLangsFromCodes(codes []string) ([]Language, bool) {
	langs := make([]Language, len(codes))
	ok := true
	for i, code := range codes {
		if lang, ok := GetLangByCode(code); ok {
			langs[i] = lang
		} else {
			ok = false
		}
	}
	return langs, ok
}
