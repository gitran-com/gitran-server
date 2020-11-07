package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wzru/gitran-server/model"
)

//GetLangs list all languages
func GetLangs(ctx *gin.Context) {
	ctx.JSON(http.StatusOK,
		model.Result{
			Success: true,
			Data: gin.H{
				"languages": model.GetLangs(),
			},
		})
}

//GetLang get a language
func GetLang(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("language_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest,
			model.Result{
				Success: false,
				Msg:     err.Error(),
				Code:    http.StatusBadRequest,
				Data:    nil,
			})
		return
	}
	lang := model.GetLangByID(id)
	if lang == nil {
		ctx.JSON(http.StatusNotFound, model.Result404)
		return
	}
	ctx.JSON(http.StatusOK,
		model.Result{
			Success: true,
			Data: gin.H{
				"language": *lang,
			},
		})
}
