package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wzru/gitran-server/service"
)

func ListProj(ctx *gin.Context) {
	service.ListProj(ctx)
}

func GetProj(ctx *gin.Context) {
	service.GetProj(ctx)
}

func CreateUserProj(ctx *gin.Context) {
	service.CreateUserProj(ctx)
}
