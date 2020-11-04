package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wzru/gitran-server/service"
)

func GetUserProjCfg(ctx *gin.Context) {
	service.GetUserProjCfg(ctx)
}

func CreateUserProjCfg(ctx *gin.Context) {
	service.CreateUserProjCfg(ctx)
}

func ListUserProjCfg(ctx *gin.Context) {
	service.ListUserProjCfg(ctx)
}

func ListUserProjBrchRule(ctx *gin.Context) {
	service.ListUserProjBrchRule(ctx)
}

func GetUserProjBrchRule(ctx *gin.Context) {
	service.GetUserProjBrchRule(ctx)
}

func CreateUserProjBrchRule(ctx *gin.Context) {
	service.CreateUserProjBrchRule(ctx)
}

func SaveUserProjCfg(ctx *gin.Context) {
	service.SaveUserProjCfg(ctx)
}

func SaveUserProjBrchRule(ctx *gin.Context) {
	service.SaveUserProjBrchRule(ctx)
}
