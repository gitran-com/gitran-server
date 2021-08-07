package service

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/config"
	"github.com/gitran-com/gitran-server/constant"
	"github.com/gitran-com/gitran-server/middleware"
	"github.com/gitran-com/gitran-server/model"
	"github.com/gitran-com/gitran-server/util"
	"github.com/google/go-github/github"
	"github.com/markbates/goth/gothic"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
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
	var (
		next string
		ok   bool
		subj = constant.SubjGithubLogin
	)
	session := sessions.Default(ctx)
	if next, ok = session.Get("next").(string); !ok {
		next = "/"
	}
	//finish github auth
	github_user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		log.Errorf("auth/github/login: %+v", err.Error())
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	//get github id
	github_id, _ := strconv.ParseInt(github_user.UserID, 10, 64)
	//check if this user has ever login
	user := model.GetUserByGithubID(github_id)
	if user == nil { //if not
		//check if this github email has ever register
		user = model.GetUserByEmail(github_user.Email)
		if user == nil { //if not, register
			user = model.NewUserFromGithub(&github_user)
			user.Create()
			next = "/login/github/new"
			subj = constant.SubjGithubFirstLogin
		} else { //else, update if diff
			if user.GithubID != github_id {
				user.GithubID = github_id
				user.Write()
			}
		}
	}
	token, expires_at, refresh_before := middleware.GenUserToken(user.Name, user.ID, subj)
	//sign token in cookie
	domain := ctx.Request.URL.Hostname()
	ctx.SetCookie("token", token, 3600, "/", domain, false, false)
	ctx.SetCookie("expires_at", fmt.Sprintf("%v", expires_at), 3600, "/", domain, false, false)
	ctx.SetCookie("refresh_before", fmt.Sprintf("%v", refresh_before), 3600, "/", domain, false, false)
	ctx.Redirect(http.StatusTemporaryRedirect, next)
}

//GitHub OAuth import callback
func AuthGithubImport(ctx *gin.Context) {
	//get next jump url
	// next, _ := gothic.GetFromSession("next", ctx.Request)
	session := sessions.Default(ctx)
	var (
		next string
		ok   bool
	)
	if next, ok = session.Get("next").(string); !ok {
		next = "/"
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
	user.GithubRepoToken = github_user.AccessToken
	user.Write()
	// fmt.Printf("github_user=%+v\ntoken=%+v\n", github_user, tk)
	ctx.Redirect(http.StatusTemporaryRedirect, next)
}

func GetGithubRepos(ctx *gin.Context) {
	user := ctx.Keys["user"].(*model.User)
	tk := user.GithubRepoToken
	if tk == "" {
		ctx.JSON(http.StatusOK, util.Result{
			Success: false,
			Msg:     "github repo unauthorized",
			Code:    constant.ErrGithubRepoUnauthorized,
		})
		return
	}
	repos, err := getGithubReposFromToken(tk)
	if err != nil {
		ctx.JSON(http.StatusOK, util.Result{
			Success: false,
			Msg:     "github repo authorize error: " + err.Error(),
			Code:    constant.ErrGithubRepoUnauthorized,
		})
		user.GithubRepoToken = ""
		user.Write()
		return
	}
	ctx.JSON(http.StatusOK, util.Result{
		Success: true,
		Data: gin.H{
			"repos": repos,
		},
	})
}

func NewGithubUser(ctx *gin.Context) {
	passwd := ctx.GetString("password")
	user := ctx.Keys["user"].(*model.User)
	user.Salt = model.GenSalt()
	user.Password = model.HashSalt(passwd, user.Salt)
	user.Write()
	ctx.JSON(http.StatusOK, util.Result{
		Success: true,
		Data:    GenUserTokenData(user, constant.SubjGithubLogin, ""),
	})
}

func getGithubReposFromToken(token string) ([]model.RepoInfo, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	repos, _, err := client.Repositories.List(ctx, "", nil)
	var repo_infos []model.RepoInfo
	if err != nil {
		log.Warnf("getGithubReposFromToken error: %+v", err.Error())
		return nil, err
	}
	for _, repo := range repos {
		repo_infos = append(repo_infos, model.RepoInfo{
			ID:        *repo.ID,
			OwnerName: *repo.Owner.Login,
			Name:      *repo.Name,
			URL:       *repo.HTMLURL,
		})
	}
	return repo_infos, nil
}
