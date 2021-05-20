package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wzru/gitran-server/model"
)

//GetUser get the user
func GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := model.GetUserByName(ctx.Param("username"))
		if user == nil {
			ctx.JSON(http.StatusUnauthorized, model.Result404)
			ctx.Abort()
			return
		}
		ctx.Set("user", user)
		ctx.Next()
	}
}
