package model

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/wzru/gitran-server/config"
	"github.com/wzru/gitran-server/constant"
)

//GetLangByID gets a language by language id
func GetLangByID(id uint) *config.Lang {
	if id < uint(len(config.Langs)) {
		if config.Langs[id].ID == id {
			return &(config.Langs[id])
		}
	}
	for _, lang := range config.Langs {
		if lang.ID == id {
			return &lang
		}
	}
	return nil
}

//GetLangByCode gets a language by language code
func GetLangByCode(code string) *config.Lang {
	for _, lang := range config.Langs {
		if lang.Code == code {
			return &lang
		}
	}
	return nil
}

//GetLangsFromString extract []Lang from string like "eng|zho-CN"
func GetLangsFromString(s string) []config.Lang {
	if s == "" {
		return nil
	}
	ss := strings.Split(s, constant.Delim)
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
