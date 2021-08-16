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
	proj_id, _ := strconv.ParseInt(ctx.Query("proj_id"), 10, 64)
	cfg := model.GetProjCfgByID(proj_id)
	cfg.Process()
}
