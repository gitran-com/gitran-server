package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/model"
)

var html = `
<html>
<body>
  Pong! From Gitran API v1.<hr>
  <a href="/api/v1/auth/github?scope=user"><Button>/api/v1/auth/github?scope=user</Button></a>
  <a href="/api/v1/auth/github?scope=repo"><Button>/api/v1/auth/github?scope=repo</Button></a>
</body>
</html>
`

func PingV1(ctx *gin.Context) {
	// ctx.String(http.StatusOK, "Pong! from Gitran API v1.")
	ctx.Header("Content-Type", "text/html")
	ctx.String(http.StatusOK, html)
}

func Test(ctx *gin.Context) {
	user_id, _ := strconv.ParseInt(ctx.Query("user_id"), 10, 64)
	proj_id, _ := strconv.ParseInt(ctx.Query("proj_id"), 10, 64)
	role, _ := strconv.ParseInt(ctx.Query("role"), 10, 64)
	model.SetUserProjRole(user_id, proj_id, model.Role(role))
	ctx.JSON(http.StatusOK, gin.H{
		"data": *model.GetUserProjRole(user_id, proj_id),
	})
}
