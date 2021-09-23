package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/model"
)

func ListProjFiles(ctx *gin.Context) {
	proj := ctx.Keys["proj"].(*model.Project)
	ctx.JSON(http.StatusOK, model.Response{
		Success: true,
		Data: gin.H{
			"proj_files": model.ListProjFiles(proj.ID),
		},
	})
}

func GetProjFile(ctx *gin.Context) {
	proj := ctx.Keys["proj"].(*model.Project)
	file_id, _ := strconv.ParseInt(ctx.Param("file_id"), 10, 64)
	pf := model.GetProjFileByID(proj.ID, file_id)
	if pf == nil {
		ctx.JSON(http.StatusNotFound, model.Resp404)
	} else {
		ctx.JSON(http.StatusOK, model.Response{
			Success: true,
			Data: gin.H{
				"proj_file": pf,
				"content":   string(pf.ReadContent()),
				"sentences": model.ListValidSents(pf.ID),
			},
		})
	}
}
