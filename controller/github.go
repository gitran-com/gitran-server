package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wzru/gitran-server/service"
)

func ListGithubRepo(ctx *gin.Context) {
	service.ListGithubRepo(ctx)
}
