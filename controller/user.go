package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/service"
)

func GetUser(ctx *gin.Context) {
	service.GetUser(ctx)
}

func UpdateUser(ctx *gin.Context) {
	service.UpdateUser(ctx)
}
