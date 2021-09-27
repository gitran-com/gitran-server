package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/service"
)

func AddLikes(ctx *gin.Context) {
	service.AddLikes(ctx)
}

func AddUnlikes(ctx *gin.Context) {
	service.AddUnlikes(ctx)
}

func DelVote(ctx *gin.Context) {
	service.DelVote(ctx)
}
