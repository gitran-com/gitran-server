package service

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/wzru/gitran-server/constant"
	"github.com/wzru/gitran-server/model"
	"github.com/wzru/gitran-server/util"
)

var (
	projMutexMap = util.NewMutexMap()
)

func ListUserProjCfg(ctx *gin.Context) {
	proj := ctx.Keys["project"].(*model.Project)
	ctx.JSON(http.StatusOK, model.Result{
		Success: true,
		Msg:     "success",
		Data: gin.H{
			"project_configs": model.ListProjCfgFromProjID(proj.ID),
		}})
	return
}

func ListUserProjBrchRule(ctx *gin.Context) {
	configID, _ := strconv.ParseUint(ctx.Param("config_id"), 10, 64)
	proj := ctx.Keys["project"].(*model.Project)
	cfg := model.GetProjCfgByID(configID)
	if cfg == nil {
		ctx.JSON(http.StatusNotFound, model.Result404)
		return
	}
	if cfg.ProjID != proj.ID {
		ctx.JSON(http.StatusNotFound, model.Result404)
		return
	}
	ctx.JSON(http.StatusOK, model.Result{
		Success: true,
		Msg:     "success",
		Data: gin.H{
			"branch_rules": model.GetBrchRuleInfosFromBrchRules(model.ListBrchRuleByCfgID(configID)),
		}})
}

func GetUserProjBrchRule(ctx *gin.Context) {
}

func CreateUserProjBrchRule(ctx *gin.Context) {
	// proj := ctx.Keys["project"].(*model.Project)
	configID, _ := strconv.ParseUint(ctx.Param("config_id"), 10, 64)
	cfg := model.GetProjCfgByID(configID)
	if cfg == nil {
		ctx.JSON(http.StatusNotFound, model.Result404)
		return
	}
	// srcBr := "refs/heads" + cfg.SrcBr
	// trnBr := "refs/heads" + cfg.TrnBr
	srcFiles := ctx.PostForm("src_files")
	trnFiles := ctx.PostForm("trn_files")
	ignFiles := ctx.PostFormArray("ign_files")
	// fmt.Printf("len=%+v, ign_files=%+v\n", len(ignFiles), ignFiles)
	rule := &model.BrchRule{
		ProjCfgID: configID,
		Status:    uint8(constant.RuleStatCreated),
		SrcFiles:  srcFiles,
		TrnFiles:  trnFiles,
		IgnFiles:  strings.Join(ignFiles, constant.Delim),
	}
	rule, err := model.NewBrchRule(rule)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     err.Error(),
			Data:    nil,
		})
		return
	}
	ctx.JSON(http.StatusCreated, model.Result{
		Success: true,
		Msg:     "分支规则创建成功",
		Data:    nil,
	})
}

//GetUserProjCfg get a user project config
func GetUserProjCfg(ctx *gin.Context) {
	//TODO
}

//CreateUserProjCfg create a new project config
func CreateUserProjCfg(ctx *gin.Context) {
	proj := ctx.Keys["project"].(*model.Project)
	srcBrName := "refs/heads/" + ctx.PostForm("src_branch")
	trnBrName := "refs/heads/" + ctx.PostForm("trn_branch")
	syncTime, _ := strconv.ParseUint(ctx.PostForm("sync_time"), 10, 64)
	pushTrans := ctx.PostForm("push_trans") == "true"
	fileName := ctx.PostForm("file_name")
	projMutexMap.Lock(proj.ID)
	defer projMutexMap.Unlock(proj.ID)
	repo, err := git.PlainOpen(proj.Path)
	wt, _ := repo.Worktree()
	//先切换到src分支
	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(srcBrName),
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     err.Error(),
			Data:    nil,
		})
		return
	}
	//然后新建trn分支
	srcHead, _ := repo.Head()
	ref := plumbing.NewHashReference(plumbing.ReferenceName(trnBrName), srcHead.Hash())
	if err := repo.Storer.SetReference(ref); err != nil {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     err.Error(),
			Data:    nil,
		})
		return
	}
	projCfg := &model.ProjCfg{
		ProjID:    proj.ID,
		SrcBr:     ctx.PostForm("src_branch"),
		TrnBr:     ctx.PostForm("trn_branch"),
		SyncTime:  syncTime,
		PushTrans: pushTrans,
		FileName:  fileName,
	}
	projCfg, err = model.NewProjCfg(projCfg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     err.Error(),
			Data:    nil,
		})
		return
	}
	ctx.JSON(http.StatusCreated, model.Result{
		Success: true,
		Msg:     "项目配置创建成功",
		Data:    nil,
	})
}

func SaveUserProjCfg(ctx *gin.Context) {

}
