package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/model"
)

//GetMe gets current user info
func GetMe(ctx *gin.Context) {
	user := ctx.Keys["user"].(*model.User)
	ctx.JSON(http.StatusOK,
		model.Response{
			Success: true,
			Data: gin.H{
				"user": user,
			},
		})
}

//EditMe edit current user info
func EditMe(ctx *gin.Context) {
	user := ctx.Keys["user"].(*model.User)
	req := model.UpdateProfileRequest{}
	if ctx.BindJSON(&req) != nil || !req.Valid() {
		ctx.JSON(http.StatusBadRequest, model.Resp400)
	} else {
		user.UpdateProfile(req.Map())
		ctx.JSON(http.StatusOK, model.Response{
			Success: true,
			Data: gin.H{
				"user": user,
			},
		})
	}
}

//GetMyProjects gets my project list
func GetMyProjects(ctx *gin.Context) {
	user := ctx.Keys["user"].(*model.User)
	projs := model.ListUserProj(user.ID)
	for i := range projs {
		projs[i].FillLangs()
	}
	ctx.JSON(http.StatusOK,
		model.Response{
			Success: true,
			Data: gin.H{
				"projs": projs,
			},
		})
}

//GetUser gets a user info
func GetUser(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	user := model.GetUserByID(id)
	if user == nil {
		ctx.JSON(http.StatusNotFound, model.Resp404)
		return
	}
	ctx.JSON(http.StatusOK,
		model.Response{
			Success: true,
			Data: gin.H{
				"user": user,
			},
		})
}
