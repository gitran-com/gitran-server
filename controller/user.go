package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/service"
)

func GetMe(ctx *gin.Context) {
	service.GetMe(ctx)
}

func EditMe(ctx *gin.Context) {
	service.EditMe(ctx)
}

func GetMyProjects(ctx *gin.Context) {
	service.GetMyProjects(ctx)
}

func GetUser(ctx *gin.Context) {
	service.GetUser(ctx)
}
