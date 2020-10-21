package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PingV1(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Pong! from Gitran API v1.")
}
