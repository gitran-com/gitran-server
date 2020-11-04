package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wzru/gitran-server/config"
	"github.com/wzru/gitran-server/model"
)

//GetLangs list all languages
func GetLangs(ctx *gin.Context) {
	ctx.JSON(http.StatusOK,
		model.Result{
			Success: true,
			Msg:     "success",
			Data: gin.H{
				"languages": config.Langs,
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
				Data:    nil,
			})
		return
	}
	lang := model.GetLangByID(uint(id))
	if lang == nil {
		ctx.JSON(http.StatusNotFound, model.Result404)
		return
	}
	ctx.JSON(http.StatusOK,
		model.Result{
			Success: true,
			Msg:     "success",
			Data: gin.H{
				"language": *lang,
			},
		})
}
