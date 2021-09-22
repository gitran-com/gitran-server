package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/model"
)

func GetProjCfg(ctx *gin.Context) {
	proj := ctx.Keys["proj"].(*model.Project)
	ctx.JSON(http.StatusOK, model.Response{
		Success: true,
		Data: gin.H{
			"proj_cfg": model.GetProjCfgByID(proj.ID),
		}})
}

//UpdateProjCfg create a new project config
func UpdateProjCfg(ctx *gin.Context) {
	var (
		req  model.UpdateProjCfgRequest
		proj = ctx.Keys["proj"].(*model.Project)
		cfg  = model.GetProjCfgByID(proj.ID)
	)
	if cfg == nil {
		ctx.JSON(http.StatusBadRequest, model.Resp400)
		return
	}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.Resp400)
		return
	}
	if !req.Valid() {
		ctx.JSON(http.StatusBadRequest, model.Resp400)
		return
	}
	cfg.UpdateProjCfg(&req)
	ctx.JSON(http.StatusOK, model.Response{
		Success: true,
		Data: gin.H{
			"proj_cfg": model.GetProjCfgByID(cfg.ID),
		},
	})
}

func PreviewProjCfg(ctx *gin.Context) {
	var (
		req  model.UpdateProjCfgRequest
		proj = ctx.Keys["proj"].(*model.Project)
		cfg  = model.GetProjCfgByID(proj.ID)
	)
	if cfg == nil {
		ctx.JSON(http.StatusBadRequest, model.Resp400)
		return
	}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.Resp400)
		return
	}
	src_files, trn_files := model.GenMultiTrnFilesFromSrcFiles(req.SrcRegs, req.TrnReg, req.IgnRegs, proj)
	ctx.JSON(http.StatusOK, model.Response{
		Success: true,
		Data: gin.H{
			"src_files": src_files,
			"trn_files": trn_files,
		},
	})
}
