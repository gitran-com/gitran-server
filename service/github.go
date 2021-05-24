package service

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	log "github.com/sirupsen/logrus"
	"github.com/wzru/gitran-server/config"
	"github.com/wzru/gitran-server/constant"
	"github.com/wzru/gitran-server/middleware"
	"github.com/wzru/gitran-server/model"
)

//GitHub OAuth
func AuthGithub(ctx *gin.Context) {
	scope := ctx.Query("scope")
	session := sessions.Default(ctx)
	session.Set("next", ctx.Request.Referer())
	session.Save()
	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "provider", "github-"+scope))
	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

//GitHub OAuth login callback
func AuthGithubLogin(ctx *gin.Context) {
	//get next jump url
	// next, _ := gothic.GetFromSession("next", ctx.Request)
	session := sessions.Default(ctx)
	next := session.Get("next").(string)
	// fmt.Printf("next=%+v\n", next)
	//finish github auth
	github_user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		log.Errorf("auth/github/login: %+v", err.Error())
		ctx.Redirect(http.StatusTemporaryRedirect, config.APP.Addr)
		return
	}
	// fmt.Printf("github_user=%+v\n", github_user)
	//get github id
	github_id, _ := strconv.ParseUint(github_user.UserID, 10, 64)
	//check if this user has ever login
	user := model.GetUserByGithubID(github_id)
	if user == nil { //if not
		//check if this github email has ever register
		user = model.GetUserByEmail(github_user.Email)
		if user == nil { //if not, register
			user, _ = model.NewUserFromGithub(&github_user)
			user, _ = model.CreateUser(user)
		} else { //else, update if diff
			if user.GithubID != github_id {
				model.UpdateUserGithubID(user, github_id)
			}
		}
	}
	//sign token in cookie
	token := middleware.GenTokenFromUser(user, "github-login")
	ctx.SetCookie("token", token, 3600, "/", ctx.Request.Host, false, false)
	if next == "" {
		ctx.Redirect(http.StatusTemporaryRedirect, config.APP.Addr)
	} else {
		ctx.Redirect(http.StatusTemporaryRedirect, next)
	}
}

//GitHub OAuth import callback
func AuthGithubImport(ctx *gin.Context) {
	//get next jump url
	// next, _ := gothic.GetFromSession("next", ctx.Request)
	session := sessions.Default(ctx)
	var next string
	var ok bool
	if next, ok = session.Get("next").(string); !ok {
		next = ""
	}
	// fmt.Printf("next=%+v\n", next)
	// finish github auth
	github_user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		log.Errorf("auth/github/login: %+v", err.Error())
		ctx.Redirect(http.StatusTemporaryRedirect, config.APP.Addr)
		return
	}
	user := ctx.Keys["user"].(*model.User)
	model.NewToken(&model.Token{
		Valid:          true,
		Source:         constant.TypeGithub,
		OwnerID:        user.ID,
		OwnerName:      github_user.Name,
		OwnerAvatarURL: github_user.AvatarURL,
		AccessToken:    github_user.AccessToken,
		Scope:          "repo",
	})
	// fmt.Printf("github_user=%+v\ntoken=%+v\n", github_user, tk)
	if next == "" {
		ctx.Redirect(http.StatusTemporaryRedirect, config.APP.Addr)
	} else {
		ctx.Redirect(http.StatusTemporaryRedirect, next)
	}
}
