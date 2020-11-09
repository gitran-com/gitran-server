package service

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
	Email     string `json:"email"`
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
	// 形成请求
	var userInfoURL = "https://api.github.com/user" // GitHub用户信息获取接口
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
	// 将响应的数据写入 userInfo 中，并返回
	var userInfo githubUserInfo
	if err := json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return &userInfo, nil
}

func AuthGithubLogin(ctx *gin.Context) {
	fmt.Printf("referer=%v\n", ctx.GetHeader("Referer"))
	state := util.RandString(32)
	stateCache.Set(state, stateValue{Referer: ctx.GetHeader("Referer")}, ghExpTime)
	ctx.JSON(http.StatusAccepted, model.Result{
		Success: true,
		Data: gin.H{
			"url": fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&scope=user&redirect_uri=%s&state=%s",
				config.Github.ClientID, config.Github.CallbackURL, state),
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
	ghToken, err := getGithubToken(fmt.Sprintf(
		"https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s&state=%s",
		config.Github.ClientID, config.Github.ClientSecret, code, state,
	))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.Result{
			Success: false,
			Code:    http.StatusInternalServerError,
			Msg:     err.Error(),
		})
		return
	}
	userInfo, err := getGithubUserInfo(ghToken)
	if err != nil {
		log.Warnf("Github get user info failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, model.Result{
			Success: false,
			Code:    http.StatusInternalServerError,
			Msg:     err.Error(),
		})
		return
	}
	fmt.Printf("user info=%+v\n", userInfo)
	user := model.GetUserByGithubID(userInfo.ID)
	//如果已注册，直接登录
	if user != nil {
		ctx.JSON(http.StatusOK, model.Result{
			Success: true,
			Data: gin.H{
				"url":   val.(stateValue).Referer,
				"token": middleware.GenTokenFromUser(user, "github-login"),
			}})
		return
	}
	ctx.JSON(http.StatusAccepted, model.Result{
		Success: true,
		Data:    nil,
	})
}

func AuthGithubImport(ctx *gin.Context) {
	fmt.Printf("referer=%v\n", ctx.GetHeader("Referer"))
	state := util.RandString(32)
	stateCache.Set(state, stateValue{Referer: ctx.GetHeader("Referer")}, ghExpTime)
	ctx.JSON(http.StatusAccepted, model.Result{
		Success: true,
		Data: gin.H{
			"url": fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&scope=repo&redirect_uri=%s&state=%s",
				config.Github.ClientID, config.Github.CallbackURL, state),
		}})
}

func AuthGithubImportCallback(ctx *gin.Context) {

}

func AuthGithubRegister(ctx *gin.Context) {
	state := ctx.Query("state")
	val, ok := stateCache.Get(state)
	if !ok {
		ctx.JSON(http.StatusNotFound, model.Result404)
		return
	}
	ctx.Set("github-user-info", val.(stateValue).GithubUserInfo)
	ctx.Set("referer", val.(stateValue).Referer)
	Register(ctx)
}

// func AuthGithubBind(ctx *gin.Context) {
// 	fmt.Printf("referer=%v\n", ctx.GetHeader("Referer"))
// 	state := util.RandString(32)
// 	stateCache.Set(state, stateValue{Referer: ctx.GetHeader("Referer")}, ghExpTime)
// 	ctx.JSON(http.StatusAccepted, model.Result{
// 		Success: true,
// 		Data: gin.H{
// 			"url": fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&scope=user&redirect_uri=%s&state=%s",
// 				config.Github.ClientID, config.Github.CallbackURL, state),
// 		}})
// }

func AuthGithubBindCallback(ctx *gin.Context) {
	user := model.GetUserByID(uint64(ctx.GetInt64("user-id")))
	if user == nil {
		ctx.JSON(http.StatusNotFound, model.Result404)
		return
	}
	if user.GithubID != 0 {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Code:    constant.ErrorGithubBindThisAccount,
			Msg:     "GitHub已绑定",
		})
		return
	}
	code := ctx.Query("code")
	state := ctx.Query("state")
	//校验state
	_, ok := stateCache.Get(state)
	if !ok {
		ctx.JSON(http.StatusNotFound, model.Result404)
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
		})
		return
	}
	userInfo, err := getGithubUserInfo(ghToken)
	if err != nil {
		log.Warnf("Github get user info failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, model.Result{
			Success: false,
			Code:    http.StatusInternalServerError,
			Msg:     err.Error(),
		})
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
		})
		return
	}
	model.UpdateUserGithubID(user, userInfo.ID)
	ctx.JSON(http.StatusOK, model.Result{
		Success: true,
		Data:    nil,
	})
}
