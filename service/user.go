package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wzru/gitran-server/model"
)

//GetUser get a user info
func GetUser(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	user := model.GetUserByID(id)
	if user == nil {
		ctx.JSON(http.StatusNotFound, model.Result404)
		return
	}
	info := model.GetUserInfoFromUser(user)
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
