package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/model"
)

//MustGetProj get a project
func MustGetProj() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uri := ctx.Param("uri")
		proj := model.GetProjByURI(uri)
		if proj == nil {
			ctx.JSON(http.StatusNotFound, model.Resp404)
			ctx.Abort()
			return
		}
		ctx.Set("proj", proj)
	}
}

//MustGetProjRole finds the project by uri
func MustGetProjRole() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			role model.Role = model.RoleNone
		)
		proj := ctx.Keys["proj"].(*model.Project)
		user := ctx.Keys["user"].(*model.User)
		//No-login user is RoleNone
		if user != nil {
			projRole := model.GetUserProjRole(user.ID, proj.ID)
			if projRole != nil {
				role = projRole.Role
			} else if proj.PublicContribute {
				role = model.RoleContributor
			} else if proj.PublicView {
				role = model.RoleViewer
			}
		}
		ctx.Set("role", role)
		ctx.Next()
	}
}
