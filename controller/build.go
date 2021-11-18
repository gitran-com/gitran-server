package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/service"
)

func NewBuild(ctx *gin.Context) {
	service.NewBuild(ctx)
}
