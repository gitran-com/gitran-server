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

func AuthGithubLogin(ctx *gin.Context) {
	service.AuthGithubLogin(ctx)
}

func AuthGithubRegister(ctx *gin.Context) {
	service.AuthGithubRegister(ctx)
}

func AuthGithubImport(ctx *gin.Context) {
	service.AuthGithubImport(ctx)
}

// func AuthGithubBind(ctx *gin.Context) {
// 	service.AuthGithubBind(ctx)
// }

func AuthGithubLoginCallback(ctx *gin.Context) {
	service.AuthGithubLoginCallback(ctx)
}

func AuthGithubImportCallback(ctx *gin.Context) {
	service.AuthGithubImportCallback(ctx)
}

func AuthGithubBindCallback(ctx *gin.Context) {
	service.AuthGithubBindCallback(ctx)
}
