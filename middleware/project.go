package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/model"
)

//MustGetProjRole finds the project by uri
func MustGetProjRole() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uri := ctx.Param("uri")
		proj := model.GetProjByURI(uri)
		if proj == nil {
			ctx.JSON(http.StatusNotFound, model.Resp404)
			ctx.Abort()
			return
		}
		ctx.Set("proj", proj)
		user := ctx.Keys["user"].(*model.User)
		if user == nil {
			ctx.Set("role", model.RoleNone)
			return
		}
		role := model.GetUserProjRole(user.ID, proj.ID)
		if role == nil {
			ctx.Set("role", model.RoleNone)
			return
		}
		ctx.Set("role", role.Role)
		ctx.Next()
	}
}
