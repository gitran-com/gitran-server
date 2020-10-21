package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wzru/gitran-server/middleware"
	"github.com/wzru/gitran-server/model"
)

func GetUser(ctx *gin.Context) {
	name := ctx.Param("name")
	user := model.GetUserByName(&name)
	info := model.GenUserInfoFromUser(user, middleware.HasPermission(ctx))
	if info != nil && info.ID != 0 {
		ctx.JSON(http.StatusOK,
			model.Result{
				Success: true,
				Msg:     "success",
				Data: gin.H{
					"user_info": *info,
				},
			})
	} else {
		ctx.JSON(http.StatusNotFound,
			model.Result{
				Success: false,
				Msg:     "not found",
				Data:    nil,
			})
	}
}

func UpdateUser(ctx *gin.Context) {
	upd := &model.UserInfo{}
	if ctx.BindJSON(upd) == nil {

	} else {
		ctx.JSON(http.StatusBadRequest,
			model.Result{
				Success: false,
				Msg:     "invalid arguments",
				Data:    nil,
			})
	}
}

func RefreshToken(ctx *gin.Context) {

}
