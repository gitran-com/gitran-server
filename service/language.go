package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/model"
)

//ListLangs list all languages
func ListLangs(ctx *gin.Context) {
	ctx.JSON(http.StatusOK,
		model.Response{
			Success: true,
			Data: gin.H{
				"langs": model.ListLangs(),
			},
		})
}

//GetLang get a language
func GetLang(ctx *gin.Context) {
	code := ctx.Param("code")
	lang, ok := model.GetLangByCode(code)
	if !ok {
		ctx.JSON(http.StatusNotFound, model.Resp404)
		return
	}
	ctx.JSON(http.StatusOK,
		model.Response{
			Success: true,
			Data: gin.H{
				"lang": lang,
			},
		})
}
