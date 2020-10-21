package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/wzru/gitran-server/config"
)

type githubToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func getToken(url string) (*githubToken, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	var client = http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var token githubToken
	if err = json.NewDecoder(res.Body).Decode(&token); err != nil {
		return nil, err
	}
	return &token, nil
}

func getGithubUserInfo(token *githubToken) (map[string]interface{}, error) {
	// 形成请求
	var userInfoUrl = "https://api.github.com/user" // github用户信息获取接口
	req, err := http.NewRequest(http.MethodGet, userInfoUrl, nil)
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
	var userInfo = make(map[string]interface{})
	if err := json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return userInfo, nil
}

func AuthGithubCallback(ctx *gin.Context) {
	code := ctx.Query("code")
	fmt.Printf("code=%v\n", code)
	url := fmt.Sprintf(
		"https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s",
		config.Github.ClientID, config.Github.ClientSecret, code,
	)
	token, err := getToken(url)
	if err != nil {
		log.Warnf("Github get token failed: %v", err)
		return
	}
	userInfo, err := getGithubUserInfo(token)
	if err != nil {
		log.Warnf("Github get user info failed: %v", err)
		return
	}
	fmt.Printf("user info=%v\n", userInfo)
}
