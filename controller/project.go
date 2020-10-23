package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wzru/gitran-server/service"
)

func GetProj(ctx *gin.Context) {
	service.GetProj(ctx)
}

func CreateUserProj(ctx *gin.Context) {
	service.CreateUserProj(ctx)
}

func CreateOrgProj(ctx *gin.Context) {
	service.CreateOrgProj(ctx)
}
