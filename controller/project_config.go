package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/service"
)

func UpdateProjCfg(ctx *gin.Context) {
	service.UpdateProjCfg(ctx)
}

func GetProjCfg(ctx *gin.Context) {
	service.GetProjCfg(ctx)
}

func PreviewProjCfg(ctx *gin.Context) {
	service.PreviewProjCfg(ctx)
}
