package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/model"
	"github.com/gitran-com/gitran-server/util"
)

//GetUser get a user info
func GetUser(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	user := model.GetUserByID(id)
	if user == nil {
		ctx.JSON(http.StatusNotFound, util.Result404)
		return
	}
	ctx.JSON(http.StatusOK,
		util.Result{
			Success: true,
			Data: gin.H{
				"user": user,
			},
		})
}

//UpdateUser update a user info
func UpdateUser(ctx *gin.Context) {
	upd := &model.User{}
	if ctx.BindJSON(upd) == nil {

	} else {
		ctx.JSON(http.StatusBadRequest,
			util.Result{
				Success: false,
				Msg:     "Invalid arguments",
				Data:    nil,
			})
	}
}
