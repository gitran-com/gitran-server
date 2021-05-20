package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/wzru/gitran-server/config"
	"github.com/wzru/gitran-server/constant"
	"github.com/wzru/gitran-server/model"
)

//AuthUserJWT verifies a token
func AuthUserJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.Request.Header.Get("Authorization")
		if len(auth) <= 0 {
			ctx.JSON(http.StatusUnauthorized, model.Result401)
			ctx.Abort()
			return
		}
		token := strings.Fields(auth)[1]
		clm, err := ParseToken(token) // 校验token
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, model.Result401)
			ctx.Abort()
			return
		}
		id, _ := strconv.ParseUint(clm.Id, 10, 64)
		user := model.GetUserByID(id)
		if user == nil {
			ctx.JSON(http.StatusUnauthorized, model.Result401)
			ctx.Abort()
			return
		}
		ctx.Set("user", user)
		ctx.Next()
	}
}

//AuthUserProjJWT verifies a jwt if can do something on a project
func AuthUserProjJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.Keys["user"].(*model.User)
		if user == nil || user.Login != ctx.Param("owner") {
			ctx.JSON(http.StatusNotFound, model.Result404)
			ctx.Abort()
			return
		}
		owner := model.GetUserByName(ctx.Param("owner"))
		if owner == nil {
			ctx.JSON(http.StatusNotFound, model.Result404)
			ctx.Abort()
			return
		}
		ctx.Set("project", proj)
		ctx.Next()
	}
}

//AuthUserGithubJWT verifies a token and its github-id
func AuthUserGithubJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.Keys["user"].(*model.User)
		if user == nil || user.GithubID == 0 {
			ctx.JSON(http.StatusUnauthorized, model.Result401)
			ctx.Abort()
			return
		}
		tk := model.GetTokenByOwnerID(user.ID, constant.TypeGithub, constant.TokenRepo)
		if tk == nil {
			ctx.JSON(http.StatusUnauthorized, model.Result401)
			ctx.Abort()
			return
		}
		ctx.Set("github-token", tk)
		ctx.Next()
	}
}

//GenTokenFromUser gen a token from User
func GenTokenFromUser(user *model.User, subj string) string {
	now := time.Now().Unix()
	claims := jwt.StandardClaims{
		Audience:  user.Name,                         // 受众
		ExpiresAt: now + int64(config.JWT.ValidTime), // 失效时间
		Id:        fmt.Sprintf("%v", user.ID),        // 编号
		IssuedAt:  now,                               // 签发时间
		Issuer:    config.APP.Name,                   // 签发人
		NotBefore: now,                               // 生效时间
		Subject:   subj,                              // 主题
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, _ := tokenClaims.SignedString([]byte(config.JWT.Secret))
	return token
}

//ParseToken parse token. Return nil claim when parse error
func ParseToken(tokenStr string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(config.JWT.Secret), nil
	})
	if token != nil {
		if claim, ok := token.Claims.(*jwt.StandardClaims); ok {
			if token.Valid {
				return claim, nil
			}
			return claim, errors.New("token is expired")
		}
	}
	return nil, err
}

//HasUserPermission check if this user has permission to user uid by checking JWT
func HasUserPermission(ctx *gin.Context, uid uint64) bool {
	auth := ctx.Request.Header.Get("Authorization")
	if len(auth) == 0 {
		return false
	}
	token := strings.Fields(auth)[1]
	clm, err := ParseToken(token)
	id, _ := strconv.ParseUint(clm.Id, 10, 64)
	return err == nil && uid == id
}
