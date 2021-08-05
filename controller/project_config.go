package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/service"
)

func CreateUserProjCfg(ctx *gin.Context) {
	service.CreateUserProjCfg(ctx)
}

func ListUserProjCfg(ctx *gin.Context) {
	service.ListUserProjCfg(ctx)
}

func ListUserProjBrchRule(ctx *gin.Context) {
	service.ListUserProjBrchRule(ctx)
}

func CreateUserProjBrchRule(ctx *gin.Context) {
	service.CreateUserProjBrchRule(ctx)
}

func SaveUserProjCfg(ctx *gin.Context) {
	service.SaveUserProjCfg(ctx)
}
