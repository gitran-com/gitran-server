package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/service"
)

func ListUserProj(ctx *gin.Context) {
	service.ListUserProj(ctx)
}

func GetProj(ctx *gin.Context) {
	service.GetProj(ctx)
}

func CreateUserProj(ctx *gin.Context) {
	service.CreateUserProj(ctx)
}

func ListProjBrch(ctx *gin.Context) {
	service.ListProjBrch(ctx)
}

func ProjExisted(ctx *gin.Context) {
	service.ProjExisted(ctx)
}

func ProjStatWS(ctx *gin.Context) {
	service.ProjStatWS(ctx)
}
