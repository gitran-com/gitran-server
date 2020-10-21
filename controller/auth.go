package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wzru/gitran-server/service"
)

func Login(ctx *gin.Context) {
	service.Login(ctx)
}

func Register(ctx *gin.Context) {
	service.Register(ctx)
}

func AuthGithubCallback(ctx *gin.Context) {
	service.AuthGithubCallback(ctx)
}
