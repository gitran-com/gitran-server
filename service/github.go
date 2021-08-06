package service

import (
	"context"
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
	//get next jump url
	// next, _ := gothic.GetFromSession("next", ctx.Request)
	session := sessions.Default(ctx)
	var (
		next string
		ok   bool
	)
	if next, ok = session.Get("next").(string); !ok {
		next = ""
	}
	//finish github auth
	github_user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		log.Errorf("auth/github/login: %+v", err.Error())
		ctx.Redirect(http.StatusTemporaryRedirect, config.APP.Addr)
		return
	}
	// fmt.Printf("github_user=%+v\n", github_user)
	//get github id
	github_id, _ := strconv.ParseInt(github_user.UserID, 10, 64)
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
	token, _, _ := middleware.GenUserToken(user.Name, user.ID, "github-login")
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
	var (
		next string
		ok   bool
	)
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

func GetGithubTokens(ctx *gin.Context) {
	user := ctx.Keys["user"].(*model.User)
	tk := model.GetValidTokensByOwnerID(user.ID, constant.TypeGithub)
	ctx.JSON(http.StatusOK, util.Result{
		Success: true,
		Data: gin.H{
			"tokens": tk,
		},
	})
}

func GetGithubRepos(ctx *gin.Context) {
	user := ctx.Keys["user"].(*model.User)
	token_id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	tk := model.GetTokenByID(token_id)
	if tk == nil || tk.OwnerID != user.ID {
		ctx.JSON(http.StatusBadRequest, util.Result404)
	} else {
		ctx.JSON(http.StatusOK, util.Result{
			Success: true,
			Data: gin.H{
				"repo_infos": getGithubReposFromToken(tk.AccessToken),
			},
		})
	}
}

func getGithubReposFromToken(token string) []model.RepoInfo {
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
	}
	for _, repo := range repos {
		repo_infos = append(repo_infos, model.RepoInfo{
			ID:        int64(*repo.ID),
			OwnerName: *repo.Owner.Login,
			Name:      *repo.Name,
			URL:       *repo.HTMLURL,
		})
	}
	return repo_infos
}
