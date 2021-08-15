package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/model"
	"github.com/gitran-com/gitran-server/util"
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
			"proj_cfg": cfg,
		},
	})
}

func processSrcFile(file string, cfg *model.ProjCfg, proj *model.Project, rule *model.FileMap) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	ext := filepath.Ext(file)
	var res []string
	switch ext {
	case ".xml":
		res = util.ProcessXML(string(data))
	}
	for _, str := range res {
		fmt.Printf("sen=%s\n", str)
	}
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
	ctx.JSON(http.StatusOK, model.Response{
		Success: true,
		Data: gin.H{
			"files": util.ListMatchFiles(proj.Path, model.GetFileMapsSrcFiles(req.FileMaps), req.IgnRegs),
		},
	})
}
