package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wzru/gitran-server/service"
)

func GetUser(ctx *gin.Context) {
	service.GetUser(ctx)
}

func UpdateUser(ctx *gin.Context) {
	service.UpdateUser(ctx)
}

func Refresh(ctx *gin.Context) {
	service.RefreshToken(ctx)
}
