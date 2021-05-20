package service

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"github.com/wzru/gitran-server/middleware"
	"github.com/wzru/gitran-server/model"
)

//GitHub OAuth
func AuthGithub(ctx *gin.Context) {
	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "provider", "github"))
	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

//GitHub OAuth callback
func AuthGithubLogin(ctx *gin.Context) {
	github_user, _ := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	github_id, _ := strconv.ParseUint(github_user.UserID, 10, 64)
	user := model.GetUserByExternID(github_id)
	if user == nil { //register
		user, _ = model.NewUserFromExtern(&github_user, model.LoginTypeGithub)
		user, _ = model.CreateUser(user)
	}
	token := middleware.GenTokenFromUser(user, "github-login")
	ctx.SetCookie("token", token, -1, "/", ctx.Request.Host, false, false)
	ctx.Redirect(http.StatusTemporaryRedirect, ctx.GetHeader("Referer"))
}
