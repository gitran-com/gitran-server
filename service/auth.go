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
	var req model.LoginRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.Resp400)
		return
	}
	user := model.GetUserByEmail(req.Email)
	if user.NoPassword || !model.CheckPass(user, req.Password) {
		ctx.JSON(http.StatusOK, model.Response{
			Success: false,
			Msg:     "email or password incorrect",
			Code:    constant.ErrEmailOrPassIncorrect,
		})
	} else {
		ctx.JSON(http.StatusOK, model.Response{
			Success: true,
			Msg:     "login successfully",
			Data:    GenUserTokenData(user, constant.SubjLogin, ctx.Request.Referer()),
		})
	}
}

//Register register new user
func Register(ctx *gin.Context) {
	var req model.RegisterRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.Resp400)
		return
	}
	user := model.GetUserByEmail(req.Email)
	if user == nil { //create new user
		salt := []byte(model.GenSalt())
		user := &model.User{
			Name:       req.Name,
			Email:      req.Email,
			Password:   model.HashSalt(req.Password, salt),
			Salt:       salt,
			IsActive:   !config.Email.Enable,
			NoPassword: false,
		}
		if err := user.Create(); err != nil {
			ctx.JSON(http.StatusOK,
				model.Response{
					Success: false,
					Msg:     err.Error(),
					Code:    constant.ErrUnknown,
				})
			return

		} else {
			ctx.JSON(http.StatusCreated,
				model.Response{
					Success: true,
					Msg:     "register successfully",
					Data:    GenUserTokenData(user, constant.SubjRegister, ctx.Request.Referer()),
				})
			return
		}
	} else {
		ctx.JSON(http.StatusOK,
			model.Response{
				Success: false,
				Msg:     "email exists",
				Code:    constant.ErrEmailExists,
			})
	}
}

//RefreshToken refresh JWT
func RefreshToken(ctx *gin.Context) {
	auth := ctx.Request.Header.Get("Authorization")
	if len(auth) == 0 {
		ctx.JSON(http.StatusOK, model.RespInvalidToken)
		ctx.Abort()
		return
	}
	token := strings.Fields(auth)[1]
	clm, _ := middleware.ParseToken(token) // 校验token
	if clm == nil {
		ctx.JSON(http.StatusOK, model.RespInvalidToken)
	} else {
		id, _ := strconv.ParseInt(clm.Id, 10, 64)
		user := model.GetUserByID(id)
		if user == nil || clm.NotBefore+config.JWT.RefreshTime < time.Now().Unix() {
			ctx.JSON(http.StatusOK, model.RespInvalidToken)
		} else {
			ctx.JSON(http.StatusOK, model.Response{
				Success: true,
				Msg:     "refresh successfully",
				Data:    GenUserTokenData(user, constant.SubjRefresh, ctx.Request.Referer()),
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
