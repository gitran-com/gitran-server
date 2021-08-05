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
)

//Login make users login
func Login(ctx *gin.Context) {
	email := ctx.PostForm("email")
	passwd := ctx.PostForm("password")
	user := model.GetUserByEmail(email)
	if model.CheckPass(user, passwd) {
		ctx.JSON(http.StatusOK, model.Result{
			Success: true,
			Msg:     "登录成功",
			Data: gin.H{
				"token": middleware.GenTokenFromUser(user, "login"),
				"url":   ctx.Request.Referer(),
			},
		})
	} else {
		ctx.JSON(http.StatusUnauthorized, model.Result{
			Success: false,
			Msg:     "用户名或密码错误",
			Code:    constant.ErrLoginOrPasswordIncorrect,
			Data:    nil,
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
				model.Result{
					Success: true,
					Msg:     "注册成功",
					Data: gin.H{
						"token": middleware.GenTokenFromUser(user, "register"),
						"url":   ctx.Request.Referer(),
					},
				})
			return
		} else {
			ctx.JSON(http.StatusBadRequest,
				model.Result{
					Success: false,
					Msg:     err.Error(),
					Data:    nil,
					Code:    -1,
				})
			return
		}
	} else {
		ctx.JSON(http.StatusBadRequest,
			model.Result{
				Success: false,
				Msg:     "邮箱不可用",
				Code:    constant.ErrEmailExists,
				Data:    nil,
			})
	}
}

//RefreshToken refresh JWT
func RefreshToken(ctx *gin.Context) {
	auth := ctx.Request.Header.Get("Authorization")
	if len(auth) == 0 {
		ctx.JSON(http.StatusUnauthorized, model.Result401)
		ctx.Abort()
		return
	}
	token := strings.Fields(auth)[1]
	clm, _ := middleware.ParseToken(token) // 校验token
	if clm == nil {
		ctx.JSON(http.StatusUnauthorized, model.Result401)
	} else {
		id, _ := strconv.Atoi(clm.Id)
		user := model.GetUserByID(int64(id))
		if user == nil {
			ctx.JSON(http.StatusNotFound, model.Result404)
		} else if clm.NotBefore+int64(config.JWT.RefreshTime) < time.Now().Unix() {
			ctx.JSON(http.StatusUnauthorized, model.Result401)
		} else {
			ctx.JSON(http.StatusOK, model.Result{
				Success: true,
				Msg:     "刷新成功",
				Data: gin.H{
					"token": middleware.GenTokenFromUser(user, "refresh"),
				},
			})
		}
	}
}
