package service

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/constant"
	"github.com/gitran-com/gitran-server/model"
	"github.com/gitran-com/gitran-server/util"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var (
	projMutexMap = util.NewMutexMap()
)

func ListUserProjCfg(ctx *gin.Context) {
	proj := ctx.Keys["project"].(*model.Project)
	ctx.JSON(http.StatusOK, util.Result{
		Success: true,
		Msg:     "success",
		Data: gin.H{
			"proj_cfgs": model.ListProjCfgByProjID(proj.ID),
		}})
}

func ListUserProjBrchRule(ctx *gin.Context) {
	configID, _ := strconv.ParseInt(ctx.Param("config_id"), 10, 64)
	proj := ctx.Keys["project"].(*model.Project)
	cfg := model.GetProjCfgByID(configID)
	if cfg == nil {
		ctx.JSON(http.StatusNotFound, util.Result404)
		return
	}
	if cfg.ProjID != proj.ID {
		ctx.JSON(http.StatusNotFound, util.Result404)
		return
	}
	ctx.JSON(http.StatusOK, util.Result{
		Success: true,
		Msg:     "success",
		Data: gin.H{
			"branch_rules": model.GetBrchRuleInfosFromBrchRules(model.ListBrchRuleByCfgID(configID)),
		}})
}

func CreateUserProjBrchRule(ctx *gin.Context) {
	configID, _ := strconv.ParseInt(ctx.Param("config_id"), 10, 64)
	cfg := model.GetProjCfgByID(configID)
	if cfg == nil {
		ctx.JSON(http.StatusNotFound, util.Result404)
		return
	}
	srcFiles := ctx.PostForm("src_files")
	trnFiles := ctx.PostForm("trn_files")
	ignFiles := ctx.PostFormArray("ign_files")
	rule := &model.BrchRule{
		ProjCfgID: configID,
		Status:    constant.RuleStatCreated,
		SrcFiles:  srcFiles,
		TrnFiles:  trnFiles,
		IgnFiles:  strings.Join(ignFiles, constant.Delim),
	}
	_, err := model.NewBrchRule(rule)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.Result{
			Success: false,
			Msg:     err.Error(),
			Data:    nil,
		})
		return
	}
	model.UpdateProjCfgChanged(cfg, true)
	ctx.JSON(http.StatusCreated, util.Result{
		Success: true,
		Msg:     "分支规则创建成功",
		Data:    nil,
	})
}

//CreateUserProjCfg create a new project config
func CreateUserProjCfg(ctx *gin.Context) {
	proj := ctx.Keys["project"].(*model.Project)
	srcBrName := "refs/heads/" + ctx.PostForm("src_branch")
	trnBrName := "refs/heads/" + ctx.PostForm("trn_branch")
	pullItv, _ := strconv.ParseUint(ctx.PostForm("pull_interval"), 10, 16)
	pushItv, _ := strconv.ParseUint(ctx.PostForm("push_interval"), 10, 16)
	fileName := ctx.PostForm("file_name")
	lk := projMutexMap.Lock(proj.ID)
	defer lk.Unlock()
	repo, _ := git.PlainOpen(proj.Path)
	wt, _ := repo.Worktree()
	//先切换到src分支
	err := wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(srcBrName),
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.Result{
			Success: false,
			Msg:     err.Error(),
			Data:    nil,
			Code:    -1,
		})
		return
	}
	//然后新建trn分支
	srcHead, _ := repo.Head()
	ref := plumbing.NewHashReference(plumbing.ReferenceName(trnBrName), srcHead.Hash())
	if err := repo.Storer.SetReference(ref); err != nil {
		ctx.JSON(http.StatusBadRequest, util.Result{
			Success: false,
			Msg:     err.Error(),
			Data:    nil,
			Code:    -1,
		})
		return
	}
	projCfg := &model.ProjCfg{
		ProjID:   proj.ID,
		SrcBr:    ctx.PostForm("src_branch"),
		TrnBr:    ctx.PostForm("trn_branch"),
		Changed:  false,
		PullItv:  uint16(pullItv),
		PushItv:  uint16(pushItv),
		FileName: fileName,
	}
	projCfg, err = model.NewProjCfg(projCfg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.Result{
			Success: false,
			Msg:     err.Error(),
			Data:    nil,
			Code:    -1,
		})
		return
	}
	ctx.JSON(http.StatusCreated, util.Result{
		Success: true,
		Msg:     "项目配置创建成功",
		Data:    nil,
	})
	if pullItv != 0 {
		pullSchd.Every(pullItv).Minutes().Do(pullGit, projCfg)
	}
	if proj.Type == constant.TypeGithub && pushItv != 0 {
		pushSchd.Every(pushItv).Minutes().Do(pushGit, projCfg)
	}
}

//SaveUserProjCfg save a project config in config file and then commit
func SaveUserProjCfg(ctx *gin.Context) {
	proj := ctx.Keys["project"].(*model.Project)
	user := ctx.Keys["user"].(*model.User)
	cfgs := model.ListProjCfgByProjID(proj.ID)
	lk := projMutexMap.Lock(proj.ID)
	defer lk.Unlock()
	repo, err := git.PlainOpen(proj.Path)
	wt, _ := repo.Worktree()
	for _, cfg := range cfgs {
		if !cfg.Changed { //if no changed, no need to commit
			continue
		}
		srcBr := "refs/heads/" + cfg.SrcBr
		err = wt.Checkout(&git.CheckoutOptions{
			Branch: plumbing.ReferenceName(srcBr),
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, util.Result{
				Success: false,
				Msg:     err.Error(),
				Data:    nil,
				Code:    constant.ErrGitChkout,
			})
			return
		}
		rules := model.ListBrchRuleByCfgID(cfg.ID)
		err = ioutil.WriteFile(proj.Path+cfg.FileName, genCfgFileByRuleInfos(model.GetBrchRuleInfosFromBrchRules(rules)), os.ModeAppend|os.ModePerm)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, util.Result{
				Success: false,
				Msg:     err.Error(),
				Data:    nil,
				Code:    -1,
			})
			return
		}
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
				ctx.JSON(http.StatusBadRequest, util.Result{
					Success: false,
					Msg:     err.Error(),
					Data:    nil,
					Code:    constant.ErrGitCommit,
				})
				return
			}
		}
		model.UpdateProjCfgChanged(&cfg, false)
		go processCfg(&cfg, proj, rules)
	}
	ctx.JSON(http.StatusCreated, util.Result{
		Success: true,
		Msg:     "项目配置更新成功",
		Data:    nil,
	})
}

func processCfg(cfg *model.ProjCfg, proj *model.Project, rules []model.BrchRule) {
	lk := projMutexMap.Lock(proj.ID)
	model.UpdateProjStatus(proj, constant.ProjStatProcessingString)
	defer lk.Unlock()
	defer model.UpdateProjStatus(proj, constant.ProjStatReady)
	repo, err := git.PlainOpen(proj.Path)
	if err != nil {
		log.Errorf("processCfg error when git.PlainOpen(): %+v", err.Error())
		return
	}
	wt, _ := repo.Worktree()
	srcBr := "refs/heads/" + cfg.SrcBr
	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(srcBr),
	})
	if err != nil {
		log.Errorf("processCfg error when Checkout(): %+v", err.Error())
		return
	}
	for _, rule := range rules {
		rule.IgnFileRegs = strings.Split(rule.IgnFiles, constant.Delim)
		files := ListMatchFiles(proj.Path, rule.SrcFiles, rule.IgnFileRegs)
		var wg sync.WaitGroup
		wg.Add(len(files))
		for _, file := range files {
			go func(file string, cfg *model.ProjCfg, proj *model.Project, rule *model.BrchRule) {
				defer wg.Done()
				abs := filepath.Join(proj.Path, file)
				fmt.Printf("abs=%+v\n", abs)
				processSrcFile(abs, cfg, proj, rule)
			}(file, cfg, proj, &rule)
		}
	}
}

func processSrcFile(file string, cfg *model.ProjCfg, proj *model.Project, rule *model.BrchRule) {
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

func ListMatchFiles(root string, pattern string, ignore []string) []string {
	files := []string{}
	// fmt.Printf("root=%s pat=%s\n", root, pattern)
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.Type().IsRegular() {
			relativePath := strings.TrimPrefix(path, root)
			// fmt.Printf("relpath=%s\n", relativePath)
			if strings.HasPrefix(relativePath, ".git/") {
				return nil
			}
			if ok, _ := filepath.Match(pattern, relativePath); !ok {
				return nil
			}
			for _, pat := range ignore {
				if ok, _ := filepath.Match(pat, relativePath); ok {
					return nil
				}
			}
			files = append(files, relativePath)
		}
		return nil
	})
	// fmt.Printf("files=%+v\n", files)
	return files
}

func genCfgFileByRuleInfos(rules []model.BrchRuleInfo) []byte {
	data, err := yaml.Marshal(map[string]interface{}{
		"rules": rules,
	})
	if err != nil {
		log.Warn(err.Error())
	}
	return data
}
