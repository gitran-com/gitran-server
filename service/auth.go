package service

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/config"
	"github.com/gitran-com/gitran-server/constant"
	"github.com/gitran-com/gitran-server/middleware"
	"github.com/gitran-com/gitran-server/model"
	"github.com/gitran-com/gitran-server/util"
)

//Login make users login
func Login(ctx *gin.Context) {
	email := ctx.PostForm("email")
	passwd := ctx.PostForm("password")
	user := model.GetUserByEmail(email)
	if model.CheckPass(user, passwd) {
		ctx.JSON(http.StatusOK, util.Result{
			Success: true,
			Msg:     "login successfully",
			Data:    GenUserTokenData(user, "login", ctx.Request.Referer()),
		})
	} else {
		ctx.JSON(http.StatusOK, util.Result{
			Success: false,
			Msg:     "email or password incorrect",
			Code:    constant.ErrEmailOrPassIncorrect,
		})
	}
}

//Register register new user
func Register(ctx *gin.Context) {
	name := ctx.PostForm("name")
	email := ctx.PostForm("email")
	passwd := ctx.PostForm("password")
	user := model.GetUserByEmail(email)
	if user == nil { //create new user
		var user *model.User
		var err error
		salt := []byte(model.GenSalt())
		user, err = model.CreateUser(&model.User{
			Name:        name,
			Email:       email,
			Password:    model.HashSalt(passwd, salt),
			Salt:        salt,
			IsActive:    !config.Email.Enable,
			LastLoginAt: time.Now(),
		})
		if err == nil {
			ctx.JSON(http.StatusCreated,
				util.Result{
					Success: true,
					Msg:     "register successfully",
					Data:    GenUserTokenData(user, "register", ctx.Request.Referer()),
				})
			return
		} else {
			ctx.JSON(http.StatusOK,
				util.Result{
					Success: false,
					Msg:     err.Error(),
					Code:    constant.ErrUnknown,
				})
			return
		}
	} else {
		ctx.JSON(http.StatusOK,
			util.Result{
				Success: false,
				Msg:     "email exists",
				Code:    constant.ErrEmailExists,
				Data:    nil,
			})
	}
}

//RefreshToken refresh JWT
func RefreshToken(ctx *gin.Context) {
	auth := ctx.Request.Header.Get("Authorization")
	if len(auth) == 0 {
		ctx.JSON(http.StatusUnauthorized, util.Result401)
		ctx.Abort()
		return
	}
	token := strings.Fields(auth)[1]
	clm, _ := middleware.ParseToken(token) // 校验token
	if clm == nil {
		ctx.JSON(http.StatusUnauthorized, util.Result401)
	} else {
		id, _ := strconv.Atoi(clm.Id)
		user := model.GetUserByID(int64(id))
		if user == nil || clm.NotBefore+int64(config.JWT.RefreshTime) < time.Now().Unix() {
			ctx.JSON(http.StatusUnauthorized, util.Result401)
		} else {
			ctx.JSON(http.StatusOK, util.Result{
				Success: true,
				Msg:     "refresh successfully",
				Data:    GenUserTokenData(user, "refresh", ctx.Request.Referer()),
			})
		}
	}
}

func GenUserTokenData(user *model.User, subj string, url string) map[string]interface{} {
	token, expired, refresh := middleware.GenUserToken(user.Name, user.ID, subj)
	dat := map[string]interface{}{
		"url":            url,
		"token":          token,
		"expires_at":     expired,
		"refresh_before": refresh,
	}
	return dat
}
