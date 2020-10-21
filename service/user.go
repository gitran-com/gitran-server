package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wzru/gitran-server/middleware"
	"github.com/wzru/gitran-server/model"
)

func GetUser(ctx *gin.Context) {
	login := ctx.Param("login")
	user := model.GetUserByLogin(login)
	if user == nil {
		ctx.JSON(http.StatusNotFound,
			model.Result{
				Success: false,
				Msg:     "Not found",
				Data:    nil,
			})
		return
	}
	info := model.GetUserInfoFromUser(user, middleware.HasUserPermission(ctx, login))
	ctx.JSON(http.StatusOK,
		model.Result{
			Success: true,
			Msg:     "success",
			Data: gin.H{
				"user_info": *info,
			},
		})
}

func UpdateUser(ctx *gin.Context) {
	upd := &model.UserInfo{}
	if ctx.BindJSON(upd) == nil {

	} else {
		ctx.JSON(http.StatusBadRequest,
			model.Result{
				Success: false,
				Msg:     "Invalid arguments",
				Data:    nil,
			})
	}
}
