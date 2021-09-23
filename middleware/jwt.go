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
	"github.com/golang-jwt/jwt"
)

//MustAuthUser verifies a token
func MustAuthUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.Request.Header.Get("Authorization")
		if len(auth) <= 0 {
			ctx.JSON(http.StatusUnauthorized, model.RespInvalidToken)
			ctx.Abort()
			return
		}
		token := strings.Fields(auth)[1]
		clm, err := ParseToken(token) // 校验token
		if err != nil || clm.Subject == constant.SubjGithubFirstLogin {
			ctx.JSON(http.StatusUnauthorized, model.RespInvalidToken)
			ctx.Abort()
			return
		}
		id, _ := strconv.ParseInt(clm.Id, 10, 64)
		user := model.GetUserByID(id)
		if user == nil {
			ctx.JSON(http.StatusUnauthorized, model.RespInvalidToken)
			ctx.Abort()
			return
		}
		ctx.Set("user", user)
		ctx.Next()
	}
}

//TryAuthUser try to verifies a token
func TryAuthUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.Request.Header.Get("Authorization")
		if len(auth) <= 0 {
			ctx.Set("user", (*model.User)(nil))
			return
		}
		token := strings.Fields(auth)[1]
		clm, err := ParseToken(token) // 校验token
		if err != nil {
			ctx.Set("user", (*model.User)(nil))
			return
		}
		id, _ := strconv.ParseInt(clm.Id, 10, 64)
		user := model.GetUserByID(id)
		ctx.Set("user", user)
	}
}

//MustAuthProjAdmin verifies a jwt if can do something on a project
func MustAuthProjAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uri := ctx.Param("uri")
		proj := model.GetProjByURI(uri)
		if proj == nil {
			ctx.JSON(http.StatusNotFound, model.Resp404)
			ctx.Abort()
			return
		}
		user := ctx.Keys["user"].(*model.User)
		role := model.GetUserProjRole(user.ID, proj.ID)
		if role == nil || role.Role > model.RoleAdmin {
			ctx.JSON(http.StatusForbidden, model.RespInvalidToken)
			ctx.Abort()
			return
		}
		ctx.Set("proj", proj)
		ctx.Next()
	}
}

//MustAuthProjCommiter verifies a jwt if can do something on a project
func MustAuthProjCommiter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uri := ctx.Param("uri")
		proj := model.GetProjByURI(uri)
		if proj == nil {
			ctx.JSON(http.StatusNotFound, model.Resp404)
			ctx.Abort()
			return
		}
		user := ctx.Keys["user"].(*model.User)
		role := model.GetUserProjRole(user.ID, proj.ID)
		if role == nil || role.Role > model.RoleCommitter {
			ctx.JSON(http.StatusForbidden, model.RespInvalidToken)
			ctx.Abort()
			return
		}
		ctx.Set("proj", proj)
		ctx.Next()
	}
}

//MustAuthProjContributor verifies a jwt if can do something on a project
func MustAuthProjContributor() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uri := ctx.Param("uri")
		proj := model.GetProjByURI(uri)
		if proj == nil {
			ctx.JSON(http.StatusNotFound, model.Resp404)
			ctx.Abort()
			return
		}
		user := ctx.Keys["user"].(*model.User)
		role := model.GetUserProjRole(user.ID, proj.ID)
		flag := true
		if role == nil {
			if !proj.PublicContribute {
				flag = false
			}
		} else if role.Role > model.RoleContributor {
			flag = false
		}
		if !flag {
			ctx.JSON(http.StatusForbidden, model.RespInvalidToken)
			ctx.Abort()
			return
		}
		ctx.Set("proj", proj)
		ctx.Next()
	}
}

//MustAuthProjViewer verifies a jwt if can do something on a project
func MustAuthProjViewer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uri := ctx.Param("uri")
		proj := model.GetProjByURI(uri)
		if proj == nil {
			ctx.JSON(http.StatusNotFound, model.Resp404)
			ctx.Abort()
			return
		}
		user := ctx.Keys["user"].(*model.User)
		role := model.GetUserProjRole(user.ID, proj.ID)
		flag := true
		if role == nil {
			if !proj.PublicView && !proj.PublicContribute {
				flag = false
			}
		} else if role.Role > model.RoleViewer {
			flag = false
		}
		if !flag {
			ctx.JSON(http.StatusForbidden, model.RespInvalidToken)
			ctx.Abort()
			return
		}
		ctx.Set("proj", proj)
		ctx.Next()
	}
}

//MustAuthNewGithubUser verifies a token
func MustAuthNewGithubUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.Request.Header.Get("Authorization")
		if len(auth) <= 0 {
			ctx.JSON(http.StatusUnauthorized, model.RespInvalidToken)
			ctx.Abort()
			return
		}
		token := strings.Fields(auth)[1]
		clm, err := ParseToken(token) // 校验token
		if err != nil || clm.Subject != constant.SubjGithubFirstLogin {
			ctx.JSON(http.StatusUnauthorized, model.RespInvalidToken)
			ctx.Abort()
			return
		}
		id, _ := strconv.ParseInt(clm.Id, 10, 64)
		user := model.GetUserByID(id)
		if user == nil {
			ctx.JSON(http.StatusUnauthorized, model.RespInvalidToken)
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
