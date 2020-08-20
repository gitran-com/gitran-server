package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/WangZhengru/gitran-be/config"
	"github.com/WangZhengru/gitran-be/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// type Claim struct {
// 	UID       uint64
// 	Name      string
// 	ExpiresAt time.Time
// }

func AuthJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		result := model.Result{
			Code: http.StatusUnauthorized,
			Msg:  "无法认证，重新登录",
			Data: nil,
		}
		auth := ctx.Request.Header.Get("Authorization")
		if len(auth) == 0 {
			ctx.Abort()
			ctx.JSON(http.StatusUnauthorized, result)
		}
		auth = strings.Fields(auth)[1]
		// 校验token
		_, err := parseToken(auth)
		if err != nil {
			ctx.Abort()
			result.Msg = "token过期; " + err.Error()
			ctx.JSON(http.StatusUnauthorized, result)
		} else {
			println("token 正确")
		}
		ctx.Next()
	}
}

func GenJWT(user *model.User, subj string) *string {
	claims := jwt.StandardClaims{
		Audience:  user.Name,                                       // 受众
		ExpiresAt: time.Now().Unix() + int64(config.JWT.ValidTime), // 失效时间
		Id:        string(user.ID),                                 // 编号
		IssuedAt:  time.Now().Unix(),                               // 签发时间
		Issuer:    "gitran",                                        // 签发人
		NotBefore: time.Now().Unix(),                               // 生效时间
		Subject:   subj,                                            // 主题
	}
	jwtSecret := []byte(config.JWT.Secret)
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, _ := tokenClaims.SignedString(jwtSecret)
	return &token
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
