package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/wzru/gitran-server/config"
	"github.com/wzru/gitran-server/model"
)

func AuthJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		result := model.Result{
			Success: false,
			Msg:     "Unauthorized",
			Data:    nil,
		}
		auth := ctx.Request.Header.Get("Authorization")
		if len(auth) == 0 {
			ctx.Abort()
			ctx.JSON(http.StatusUnauthorized, result)
		}
		auth = strings.Fields(auth)[1]
		_, err := parseToken(auth) // 校验token
		if err != nil {
			ctx.Abort()
			ctx.JSON(http.StatusUnauthorized, result)
		}
		ctx.Next()
	}
}

func GenJWT(user *model.User, subj string) *string {
	now := time.Now().Unix()
	claims := jwt.StandardClaims{
		Audience:  user.Name,                         // 受众
		ExpiresAt: now + int64(config.JWT.ValidTime), // 失效时间
		Id:        string(user.ID),                   // 编号
		IssuedAt:  now,                               // 签发时间
		Issuer:    config.APP.Name,                   // 签发人
		NotBefore: now,                               // 生效时间
		Subject:   subj,                              // 主题
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, _ := tokenClaims.SignedString([]byte(config.JWT.Secret))
	return &token
}

func HasPermission(ctx *gin.Context) bool {
	auth := ctx.Request.Header.Get("Authorization")
	if len(auth) == 0 {
		return false
	}
	auth = strings.Fields(auth)[1]
	_, err := parseToken(auth)
	return err == nil
}

func parseToken(token string) (*jwt.StandardClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(config.C.JWT.Secret), nil
	})
	if err == nil && jwtToken != nil {
		if claim, ok := jwtToken.Claims.(*jwt.StandardClaims); ok && jwtToken.Valid {
			return claim, nil
		}
	}
	return nil, err
}
