package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/model"
	"github.com/gitran-com/gitran-server/util"
)

//GetUser gets a user info
func GetMe(ctx *gin.Context) {
	user := ctx.Keys["user"].(*model.User)
	ctx.JSON(http.StatusOK,
		util.Response{
			Success: true,
			Data: gin.H{
				"user": user,
			},
		})
}

//GetUser gets a user info
func GetUser(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	user := model.GetUserByID(id)
	if user == nil {
		ctx.JSON(http.StatusNotFound, util.Resp404)
		return
	}
	ctx.JSON(http.StatusOK,
		util.Response{
			Success: true,
			Data: gin.H{
				"user": user,
			},
		})
}

//UpdateUser updates a user info
func UpdateUser(ctx *gin.Context) {
	upd := &model.User{}
	if ctx.BindJSON(upd) == nil {

	} else {
		ctx.JSON(http.StatusBadRequest,
			util.Response{
				Success: false,
				Msg:     "invalid arguments",
				Data:    nil,
			})
	}
}
