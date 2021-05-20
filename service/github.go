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

//ListGithubRepo list all github repo from auth user
func ListGithubRepo(ctx *gin.Context) {
	tk := ctx.Keys["github-token"].(*model.Token)
	// getGithubUserInfo(genGithubToken(tk))
	repos, err := getGithubRepo(genGithubToken(tk))
	if err != nil {
		log.Warnf("list github repo ERROR : %v", err.Error())
		ctx.JSON(http.StatusInternalServerError, model.Result{
			Success: false,
			Msg:     err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, model.Result{
		Success: true,
		Data: gin.H{
			"repo_infos": repos,
		}})
}

func genGithubToken(tk *model.Token) *githubToken {
	return &githubToken{
		AccessToken: tk.AccessToken,
		TokenType:   tk.TokenType,
		Scope:       tk.Scope,
	}
}

func getGithubRepoByTokenID(rid uint64, tk *model.Token) *githubRepoInfo {
	url := fmt.Sprintf("https://api.github.com/repositories/%v", rid) // GitHub仓库信息获取接口
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", "token "+tk.AccessToken)
	var client = http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil
	}
	var repoInfo githubRepoInfo
	if err := json.NewDecoder(res.Body).Decode(&repoInfo); err != nil {
		return nil
	}
	return &repoInfo
}

func getGithubRepoBrch(rid uint64, tk *model.Token) []repoBrch {
	url := fmt.Sprintf("https://api.github.com/repositories/%v/branches", rid)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Warnf("new request : %+v", err.Error())
		return nil
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", "token "+tk.AccessToken)
	var client = http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Warnf("send request : %+v", err.Error())
		return nil
	}
	var repoInfo []repoBrch
	if err := json.NewDecoder(res.Body).Decode(&repoInfo); err != nil {
		log.Warnf("json decode : %+v", err.Error())
		return nil
	}
	return repoInfo
}

//ListGithubRepoBrch list github repo branches
func ListGithubRepoBrch(ctx *gin.Context) {
	proj := ctx.Keys["project"].(*model.Project)
	tk := model.GetTokenByOwnerID(proj.OwnerID, constant.TypeGithub, constant.TokenRepo)
	if tk == nil {
		ctx.JSON(http.StatusNotFound, model.Result404)
		return
	}
	ctx.JSON(http.StatusOK, model.Result{
		Success: true,
		Data: gin.H{
			"branches": getGithubRepoBrch(proj.RepoID, tk),
		}})
	return
}

//ListGitRepoBrch list local git repo branches
func ListGitRepoBrch(ctx *gin.Context) {
	proj := ctx.Keys["project"].(*model.Project)
	repo, err := git.PlainOpen(proj.Path)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.Result{
			Success: false,
			Msg:     err.Error(),
		})
		return
	}
	it, err := repo.Branches()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.Result{
			Success: false,
			Msg:     err.Error(),
		})
		return
	}
	var brs []repoBrch
	it.ForEach(func(r *plumbing.Reference) error {
		brs = append(brs, repoBrch{
			Name: r.Name().Short(),
		})
		return nil
	})
	ctx.JSON(http.StatusOK, model.Result{
		Success: true,
		Data: gin.H{
			"branches": brs,
		}})
}
