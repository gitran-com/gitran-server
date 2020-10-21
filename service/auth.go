package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wzru/gitran-server/middleware"
	"github.com/wzru/gitran-server/model"
)

func Login(ctx *gin.Context) {
	login := ctx.PostForm("login")
	passwd := ctx.PostForm("password")
	user := model.GetUserByNameEmail(&login, &login)
	if user != nil && model.CheckPassword(&passwd, user) {
		ctx.JSON(http.StatusOK, model.Result{
			Success: true,
			Msg:     "登录成功",
			Data: gin.H{
				"token": *middleware.GenJWT(user, "login"),
			},
		})
		return
	}
	ctx.JSON(http.StatusUnauthorized, model.Result{
		Success: false,
		Msg:     "用户名或密码错误",
		Data:    nil,
	})
}

func Register(ctx *gin.Context) {
	login := ctx.PostForm("login")
	email := ctx.PostForm("email")
	passwd := ctx.PostForm("password")
	user := model.GetUserByNameEmail(&login, &email)
	if user == nil || user.ID == 0 {
		salt := []byte(model.GenSalt())
		user, err := model.NewUser(&model.User{
			Name:     login,
			Email:    email,
			Password: model.HashSalt(&passwd, salt),
			Salt:     salt,
		})
		if err == nil {
			ctx.JSON(http.StatusOK,
				model.Result{
					Success: true,
					Msg:     "注册成功",
					Data: gin.H{
						"token": *middleware.GenJWT(user, "register"),
					},
				})
			return
		} else {
			ctx.JSON(http.StatusOK,
				model.Result{
					Success: false,
					Msg:     err.Error(),
					Data:    nil,
				})
			return
		}
	}
	msg := ""
	if user.Name == login {
		msg = "用户名不可用"
	} else {
		msg = "邮箱不可用"
	}
	ctx.JSON(http.StatusOK,
		model.Result{
			Success: false,
			Msg:     msg,
			Data:    nil,
		})
}
