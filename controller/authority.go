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

func Refresh(ctx *gin.Context) {
	service.RefreshToken(ctx)
}

func AuthGithub(ctx *gin.Context) {
	service.AuthGithub(ctx)
}

func AuthGithubLogin(ctx *gin.Context) {
	service.AuthGithubLogin(ctx)
}

func AuthGithubImport(ctx *gin.Context) {
	service.AuthGithubImport(ctx)
}
