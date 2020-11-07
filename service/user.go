package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wzru/gitran-server/middleware"
	"github.com/wzru/gitran-server/model"
)

//GetUser get a user info
func GetUser(ctx *gin.Context) {
	login := ctx.Param("username")
	user := model.GetUserByName(login)
	if user == nil {
		ctx.JSON(http.StatusNotFound, model.Result404)
		return
	}
	info := model.GetUserInfoFromUser(user, middleware.HasUserPermission(ctx, user.ID))
	ctx.JSON(http.StatusOK,
		model.Result{
			Success: true,
			Data: gin.H{
				"user_info": *info,
			},
		})
}

//UpdateUser update a user info
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
