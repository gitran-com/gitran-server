package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	log "github.com/sirupsen/logrus"
	"github.com/wzru/gitran-server/constant"
	"github.com/wzru/gitran-server/model"
	"github.com/wzru/gitran-server/util"
	"gopkg.in/yaml.v2"
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
	model.UpdateProjCfgChanged(cfg, true)
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
		Changed:   false,
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

//SaveUserProjCfg save a project config
func SaveUserProjCfg(ctx *gin.Context) {
	proj := ctx.Keys["project"].(*model.Project)
	user := ctx.Keys["user"].(*model.User)
	cfgs := model.ListProjCfgFromProjID(proj.ID)
	projMutexMap.Lock(proj.ID)
	defer projMutexMap.Unlock(proj.ID)
	repo, err := git.PlainOpen(proj.Path)
	wt, _ := repo.Worktree()
	for _, cfg := range cfgs {
		if cfg.Changed == false {
			continue
		}
		srcBr := "refs/heads/" + cfg.SrcBr
		err = wt.Checkout(&git.CheckoutOptions{
			Branch: plumbing.ReferenceName(srcBr),
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, model.Result{
				Success: false,
				Msg:     err.Error(),
				Data:    nil,
			})
			return
		}
		rules := model.ListBrchRuleByCfgID(cfg.ID)
		err = ioutil.WriteFile(proj.Path+cfg.FileName, genCfgFileFromRuleInfos(model.GetBrchRuleInfosFromBrchRules(rules)), os.ModeAppend|os.ModePerm)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, model.Result{
				Success: false,
				Msg:     err.Error(),
				Data:    nil,
			})
			return
		}
		go model.UpdateProjCfgChanged(&cfg, false)
		wt.Add(cfg.FileName)
		status, _ := wt.Status()
		if len(status) > 0 {
			var msg string
			if status[cfg.FileName].Staging == git.Added {
				msg = fmt.Sprintf("feat: add Gitran config file '%s'", cfg.FileName)
			} else {
				msg = fmt.Sprintf("feat: update Gitran config file '%s'", cfg.FileName)
			}
			_, err := wt.Commit(msg, &git.CommitOptions{
				Author: &object.Signature{
					Name:  user.Name,
					Email: user.Email,
					When:  time.Now(),
				}})
			if err != nil {
				log.Warn(err.Error())
				ctx.JSON(http.StatusBadRequest, model.Result{
					Success: false,
					Msg:     err.Error(),
					Data:    nil,
				})
				return
			}
		}
	}
	ctx.JSON(http.StatusCreated, model.Result{
		Success: true,
		Msg:     "项目配置更新成功",
		Data:    nil,
	})
}

func SaveUserProjBrchRule(ctx *gin.Context) {
	//NO
}

func genCfgFileFromRuleInfos(rules []model.BrchRuleInfo) []byte {
	data, err := yaml.Marshal(map[string]interface{}{
		"rules": rules,
	})
	if err != nil {
		log.Warn(err.Error())
	}
	return data
}
