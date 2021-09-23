package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/service"
)

func ListProjFiles(ctx *gin.Context) {
	service.ListProjFiles(ctx)
}

func GetProjFile(ctx *gin.Context) {
	service.GetProjFile(ctx)
}
