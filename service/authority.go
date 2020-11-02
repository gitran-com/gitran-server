package service

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wzru/gitran-server/middleware"
	"github.com/wzru/gitran-server/model"
)

//Login make users login
func Login(ctx *gin.Context) {
	login := ctx.PostForm("login")
	passwd := ctx.PostForm("password")
	user := model.GetUserByUnameEmail(login, login)
	// fmt.Printf("user=%v\n", user)
	if user != nil && model.CheckPasswordCorrect(passwd, user) {
		ctx.JSON(http.StatusOK, model.Result{
			Success: true,
			Msg:     "登录成功",
			Data: gin.H{
				"token": middleware.GenTokenFromUser(user, "login"),
			},
		})
	} else {
		ctx.JSON(http.StatusUnauthorized, model.Result{
			Success: false,
			Msg:     "用户名或密码错误",
			Data:    nil,
		})
	}
}

//Register register new user
func Register(ctx *gin.Context) {
	login := ctx.PostForm("login")
	email := ctx.PostForm("email")
	passwd := ctx.PostForm("password")
	user := model.GetUserByUnameEmail(login, email)
	if user == nil {
		salt := []byte(model.GenSalt())
		user, err := model.NewUser(&model.User{
			Login:    login,
			Name:     login,
			Email:    email,
			Password: model.HashSalt(passwd, salt),
			Salt:     salt,
		})
		if err == nil {
			ctx.JSON(http.StatusCreated,
				model.Result{
					Success: true,
					Msg:     "注册成功",
					Data: gin.H{
						"token": middleware.GenTokenFromUser(user, "register"),
					},
				})
			return
		} else {
			ctx.JSON(http.StatusBadRequest,
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
	ctx.JSON(http.StatusBadRequest,
		model.Result{
			Success: false,
			Msg:     msg,
			Data:    nil,
		})
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
		user := model.GetUserByID(uint64(id))
		ctx.JSON(http.StatusOK, model.Result{
			Success: true,
			Msg:     "刷新成功",
			Data: gin.H{
				"token": middleware.GenTokenFromUser(user, "refresh"),
			},
		})
	}
}
