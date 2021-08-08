package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/service"
)

//ListLangs list all languages
func ListLangs(ctx *gin.Context) {
	service.ListLangs(ctx)
}

//GetLang get a language
func GetLang(ctx *gin.Context) {
	service.GetLang(ctx)
}
