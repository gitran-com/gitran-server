package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"github.com/wzru/gitran-server/config"
	"github.com/wzru/gitran-server/constant"
	"github.com/wzru/gitran-server/middleware"
	"github.com/wzru/gitran-server/model"
	"github.com/wzru/gitran-server/util"
)

var (
	ghExpTime  = 10 * time.Minute
	stateCache = cache.New(ghExpTime, 2*ghExpTime)
)

type githubUserInfo struct {
	ID        uint64 `json:"id"`
	Login     string `json:"login"`
	Type      string `json:"type"`
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
	Email     string `json:"email"`
}

type githubRepoInfo struct {
	ID      uint64         `json:"id"`
	Name    string         `json:"name"`
	URL     string         `json:"url"`
	Private bool           `json:"private"`
	Owner   githubUserInfo `json:"owner"`
}

type githubToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type stateValue struct {
	Referer        string
	GithubUserInfo githubUserInfo
}

func genState() string {
	return util.RandString(32)
}

func getGithubToken(url string) (*githubToken, error) {
	req, err := http.NewRequest(http.MethodPost, url, nil)
	req.Header.Set("Accept", "application/json")
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	//解析GitHub token
	ghToken := githubToken{}
	if err = json.NewDecoder(res.Body).Decode(&ghToken); err != nil {
		return nil, err
	}
	return &ghToken, nil
}

func getGithubUserInfo(token *githubToken) (*githubUserInfo, error) {
	var userRepoURL = "https://api.github.com/user" // GitHub用户信息获取接口
	req, err := http.NewRequest(http.MethodGet, userRepoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", "token "+token.AccessToken)
	var client = http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// 将响应的数据写入 userInfo 中，并返回
	var userInfo githubUserInfo
	if err := json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return &userInfo, nil
}

func getGithubRepo(token *githubToken) ([]githubRepoInfo, error) {
	var userInfoURL = "https://api.github.com/user/repos" // GitHub仓库信息获取接口
	req, err := http.NewRequest(http.MethodGet, userInfoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", "token "+token.AccessToken)
	var client = http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var repoInfo []githubRepoInfo
	if err := json.NewDecoder(res.Body).Decode(&repoInfo); err != nil {
		return nil, err
	}
	return repoInfo, nil
}

func storeGhToken(ghToken *githubToken, oid uint64) {
	if ghToken == nil {
		return
	}
	for _, scope := range strings.Split(ghToken.Scope, constant.GithubTokenScopeDelim) {
		if scope == "repo" {
			model.NewToken(
				&model.Token{
					Source:      constant.TypeGithub,
					OwnerID:     oid,
					AccessToken: ghToken.AccessToken,
					TokenType:   ghToken.TokenType,
					Scope:       scope,
				})
		}
	}
}

func AuthGithubLogin(ctx *gin.Context) {
	fmt.Printf("referer=%v\n", ctx.GetHeader("Referer"))
	state := genState()
	stateCache.Set(state, stateValue{Referer: ctx.GetHeader("Referer")}, ghExpTime)
	ctx.JSON(http.StatusAccepted, model.Result{
		Success: true,
		Data: gin.H{
			"url": fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&scope=user&redirect_uri=%s&state=%s",
				config.Github.ClientID, config.Github.CallbackURL+"/login", state),
		}})
}

func AuthGithubLoginCallback(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")
	//校验state
	val, ok := stateCache.Get(state)
	if !ok {
		ctx.JSON(http.StatusNotFound, model.Result404)
		return
	}
	ref := val.(stateValue).Referer
	ghToken, err := getGithubToken(fmt.Sprintf(
		"https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s&state=%s",
		config.Github.ClientID, config.Github.ClientSecret, code, state,
	))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.Result{
			Success: false,
			Code:    http.StatusInternalServerError,
			Msg:     err.Error(),
			Data: gin.H{
				"url": ref,
			}})
		return
	}
	userInfo, err := getGithubUserInfo(ghToken)
	if err != nil {
		log.Warnf("Github get user info failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, model.Result{
			Success: false,
			Code:    http.StatusInternalServerError,
			Msg:     err.Error(),
			Data: gin.H{
				"url": ref,
			}})
		return
	}
	fmt.Printf("user info=%+v\n", userInfo)
	user := model.GetUserByGithubID(userInfo.ID)
	//如果已注册，直接登录
	if user != nil {
		ctx.JSON(http.StatusOK, model.Result{
			Success: true,
			Data: gin.H{
				"url":   ref,
				"token": middleware.GenTokenFromUser(user, "github-login"),
			}})
		return
	}
	//否则要求跳转注册页面
	ctx.JSON(http.StatusAccepted, model.Result{
		Success: true,
		Data:    nil,
	})
}

func AuthGithubImport(ctx *gin.Context) {
	fmt.Printf("referer=%v\n", ctx.GetHeader("Referer"))
	state := genState()
	stateCache.Set(state, stateValue{Referer: ctx.GetHeader("Referer")}, ghExpTime)
	ctx.JSON(http.StatusAccepted, model.Result{
		Success: true,
		Data: gin.H{
			"url": fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&scope=repo&redirect_uri=%s&state=%s",
				config.Github.ClientID, config.Github.CallbackURL+"import", state),
		}})
}

func AuthGithubImportCallback(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")
	//校验state
	val, ok := stateCache.Get(state)
	if !ok {
		ctx.JSON(http.StatusNotFound, model.Result404)
		return
	}
	ref := val.(stateValue).Referer
	ghToken, err := getGithubToken(fmt.Sprintf(
		"https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s&state=%s",
		config.Github.ClientID, config.Github.ClientSecret, code, state,
	))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.Result{
			Success: false,
			Code:    http.StatusInternalServerError,
			Msg:     err.Error(),
			Data: gin.H{
				"url": ref,
			}})
		return
	}
	user := ctx.Keys["user"].(*model.User)
	uid := user.ID
	storeGhToken(ghToken, uid)
	ctx.JSON(http.StatusOK, model.Result{
		Success: true,
		Data:    nil,
	})
}

func AuthGithubRegister(ctx *gin.Context) {
	state := ctx.Query("state")
	val, ok := stateCache.Get(state)
	if !ok {
		ctx.JSON(http.StatusNotFound, model.Result404)
		return
	}
	sv := val.(stateValue)
	ctx.Set("github-user-info", &(sv.GithubUserInfo))
	ctx.Set("referer", sv.Referer)
	Register(ctx)
}

func AuthGithubBind(ctx *gin.Context) {
	fmt.Printf("referer=%v\n", ctx.GetHeader("Referer"))
	state := genState()
	stateCache.Set(state, stateValue{Referer: ctx.GetHeader("Referer")}, ghExpTime)
	ctx.JSON(http.StatusAccepted, model.Result{
		Success: true,
		Data: gin.H{
			"url": fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&scope=user&redirect_uri=%s&state=%s",
				config.Github.ClientID, config.Github.CallbackURL+"/bind", state),
		}})
}

func AuthGithubBindCallback(ctx *gin.Context) {
	user := ctx.Keys["user"].(*model.User)
	if user == nil {
		ctx.JSON(http.StatusNotFound, model.Result404)
		return
	}
	code := ctx.Query("code")
	state := ctx.Query("state")
	//校验state
	val, ok := stateCache.Get(state)
	if !ok {
		ctx.JSON(http.StatusNotFound, model.Result404)
		return
	}
	ref := val.(stateValue).Referer
	if user.GithubID != 0 {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Code:    constant.ErrorGithubBindThisAccount,
			Msg:     "GitHub已绑定",
			Data: gin.H{
				"url": ref,
			}})
		return
	}
	ghToken, err := getGithubToken(fmt.Sprintf(
		"https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s&state=%s",
		config.Github.ClientID, config.Github.ClientSecret, code, state,
	))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.Result{
			Success: false,
			Code:    http.StatusInternalServerError,
			Msg:     err.Error(),
			Data: gin.H{
				"url": ref,
			}})
		return
	}
	userInfo, err := getGithubUserInfo(ghToken)
	if err != nil {
		log.Warnf("Github get user info failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, model.Result{
			Success: false,
			Code:    http.StatusInternalServerError,
			Msg:     err.Error(),
			Data: gin.H{
				"url": ref,
			}})
		return
	}
	fmt.Printf("user info=%+v\n", userInfo)
	ghUser := model.GetUserByGithubID(userInfo.ID)
	//如果已绑定了别的账号
	if ghUser != nil {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Code:    constant.ErrorGithubBindOtherAccount,
			Msg:     "GitHub账号已绑定其他账号",
			Data: gin.H{
				"url": ref,
			}})
		return
	}
	model.UpdateUserGithubID(user, userInfo.ID)
	ctx.JSON(http.StatusOK, model.Result{
		Success: true,
		Data: gin.H{
			"url": ref,
		}})
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
