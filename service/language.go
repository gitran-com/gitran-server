package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/model"
	"github.com/gitran-com/gitran-server/util"
)

//GetLangs list all languages
func GetLangs(ctx *gin.Context) {
	ctx.JSON(http.StatusOK,
		util.Result{
			Success: true,
			Data: gin.H{
				"languages": model.GetLangs(),
			},
		})
}

//GetLang get a language
func GetLang(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest,
			util.Result{
				Success: false,
				Msg:     err.Error(),
				Code:    http.StatusBadRequest,
				Data:    nil,
			})
		return
	}
	lang := model.GetLangByID(id)
	if lang == nil {
		ctx.JSON(http.StatusNotFound, util.Result404)
		return
	}
	ctx.JSON(http.StatusOK,
		util.Result{
			Success: true,
			Data: gin.H{
				"language": *lang,
			},
		})
}
