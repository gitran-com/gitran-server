package service

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wzru/gitran-server/config"
	"github.com/wzru/gitran-server/constant"
	"github.com/wzru/gitran-server/model"
)

func checkLangCode(code string) bool {
	for _, lang := range config.Langs {
		if lang.Code == code {
			return true
		}
	}
	return false
}

func checkLangCodes(codes []string) bool {
	for _, code := range codes {
		ok := checkLangCode(code)
		if !ok {
			return false
		}
	}
	return true
}

func createGitProj(ctx *gin.Context) error {
	return nil
}

func CreateProj(ctx *gin.Context) {
	name := ctx.PostForm("name")
	tp := ctx.PostForm("type")
	desc := ctx.PostForm("description")
	priv := ctx.PostForm("private") == "1"
	src := ctx.PostForm("source_languages")
	tgt := ctx.PostForm("target_languages")
	// url := ctx.PostForm("url")
	// branch := ctx.PostForm("branch")
	if !model.CheckProjName(name) {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     "名字不合法",
			Data:    nil,
		})
		return
	}
	if src == "" {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     "源语言不能为空",
			Data:    nil,
		})
		return
	}
	srcCodes := strings.Split(src, constant.Delimiter)
	tgtCodes := strings.Split(tgt, constant.Delimiter)
	if !checkLangCodes(srcCodes) {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     "源语言不合法",
			Data:    nil,
		})
		return
	}
	if !checkLangCodes(tgtCodes) {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     "目标语言不合法",
			Data:    nil,
		})
		return
	}
	proj := &model.Project{
		Name:     name,
		Desc:     desc,
		Private:  priv,
		Type:     tp,
		Creator:  0,
		SrcLangs: src,
		TgtLangs: tgt,
		Path:     "",
	}
	if tp == constant.ProjGen {
		model.NewProject(proj)
	} else if tp == constant.ProjGit {
		if err := createGitProj(ctx); err == nil {

		}
	} else {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     "项目类型不合法",
			Data:    nil,
		})
		return
	}
}
