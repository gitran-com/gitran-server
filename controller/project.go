package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wzru/gitran-server/service"
)

func CreateProj(ctx *gin.Context) {
	service.CreateProj(ctx)
}
