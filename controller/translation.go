package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/service"
)

func ListSentTrans(ctx *gin.Context) {
	service.ListSentTrans(ctx)
}

func PostTran(ctx *gin.Context) {
	service.PostTran(ctx)
}