package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/config"
	"github.com/gitran-com/gitran-server/constant"
	"github.com/gitran-com/gitran-server/model"
	"github.com/gitran-com/gitran-server/util"
	"github.com/golang-jwt/jwt"
)

//AuthUserJWT verifies a token
func AuthUserJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.Request.Header.Get("Authorization")
		if len(auth) <= 0 {
			ctx.JSON(http.StatusUnauthorized, util.RespInvalidToken)
			ctx.Abort()
			return
		}
		token := strings.Fields(auth)[1]
		clm, err := ParseToken(token) // 校验token
		if err != nil || clm.Subject == constant.SubjGithubFirstLogin {
			ctx.JSON(http.StatusUnauthorized, util.RespInvalidToken)
			ctx.Abort()
			return
		}
		id, _ := strconv.ParseInt(clm.Id, 10, 64)
		user := model.GetUserByID(id)
		if user == nil {
			ctx.JSON(http.StatusUnauthorized, util.RespInvalidToken)
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
		id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
		proj := model.GetProjByID(id)
		if proj == nil {
			ctx.JSON(http.StatusNotFound, util.Resp404)
			ctx.Abort()
			return
		}
		user := ctx.Keys["user"].(*model.User)
		if proj.OwnerID != user.ID {
			ctx.JSON(http.StatusNotFound, util.RespInvalidToken)
			ctx.Abort()
			return
		}
		ctx.Set("project", proj)
		ctx.Next()
	}
}

//AuthNewGithubUserJWT verifies a token
func AuthNewGithubUserJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.Request.Header.Get("Authorization")
		if len(auth) <= 0 {
			ctx.JSON(http.StatusUnauthorized, util.RespInvalidToken)
			ctx.Abort()
			return
		}
		token := strings.Fields(auth)[1]
		clm, err := ParseToken(token) // 校验token
		if err != nil || clm.Subject != constant.SubjGithubFirstLogin {
			ctx.JSON(http.StatusUnauthorized, util.RespInvalidToken)
			ctx.Abort()
			return
		}
		id, _ := strconv.ParseInt(clm.Id, 10, 64)
		user := model.GetUserByID(id)
		if user == nil {
			ctx.JSON(http.StatusUnauthorized, util.RespInvalidToken)
			ctx.Abort()
			return
		}
		ctx.Set("user", user)
		ctx.Next()
	}
}

//GenUserToken gen a token from User
func GenUserToken(audience string, id int64, subj string) (string, int64, int64) {
	now := time.Now().Unix()
	claims := jwt.StandardClaims{
		Audience:  audience,                   // 受众
		ExpiresAt: now + config.JWT.ValidTime, // 失效时间
		Id:        fmt.Sprintf("%v", id),      // 编号
		IssuedAt:  now,                        // 签发时间
		Issuer:    config.APP.Name,            // 签发人
		NotBefore: now,                        // 生效时间
		Subject:   subj,                       // 主题
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, _ := tokenClaims.SignedString([]byte(config.JWT.Secret))
	return token, claims.ExpiresAt, claims.NotBefore + config.JWT.RefreshTime
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
