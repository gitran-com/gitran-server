package service

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/constant"
	"github.com/gitran-com/gitran-server/model"
	"github.com/gitran-com/gitran-server/util"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	log "github.com/sirupsen/logrus"
)

var (
	projMutexMap = util.NewMutexMap()
)

func GetProjCfg(ctx *gin.Context) {
	proj := ctx.Keys["proj"].(*model.Project)
	ctx.JSON(http.StatusOK, model.Response{
		Success: true,
		Data: gin.H{
			"proj_cfg": model.GetCfgByProjID(proj.ID),
		}})
}

//UpdateProjCfg create a new project config
func UpdateProjCfg(ctx *gin.Context) {
	var (
		req  model.UpdateProjCfgRequest
		proj = ctx.Keys["proj"].(*model.Project)
		cfg  = model.GetCfgByProjID(proj.ID)
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
	cfg.UpdateProjCfg(req.Map())
	ctx.JSON(http.StatusOK, model.Response{
		Success: true,
		Data: gin.H{
			"proj_cfg": cfg,
		},
	})
}

func processCfg(cfg *model.ProjCfg, proj *model.Project, rules []model.FileMap) {
	lk := projMutexMap.Lock(proj.ID)
	proj.UpdateStatus(constant.ProjStatProcessingString)
	defer lk.Unlock()
	defer proj.UpdateStatus(constant.ProjStatReady)
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
		files := ListMatchFiles(proj.Path, rule.SrcFileReg, cfg.IgnRegs)
		var wg sync.WaitGroup
		wg.Add(len(files))
		for _, file := range files {
			go func(file string, cfg *model.ProjCfg, proj *model.Project, rule *model.FileMap) {
				defer wg.Done()
				abs := filepath.Join(proj.Path, file)
				fmt.Printf("abs=%+v\n", abs)
				processSrcFile(abs, cfg, proj, rule)
			}(file, cfg, proj, &rule)
		}
	}
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
