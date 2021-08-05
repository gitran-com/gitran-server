package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/service"
)

//GetLangs list all languages
func GetLangs(ctx *gin.Context) {
	service.GetLangs(ctx)
}

//GetLang get a language
func GetLang(ctx *gin.Context) {
	service.GetLang(ctx)
}
