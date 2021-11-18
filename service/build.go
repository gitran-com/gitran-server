package service

import (
	"bytes"
	"net/http"
	"os"
	"path"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/config"
	"github.com/gitran-com/gitran-server/constant"
	"github.com/gitran-com/gitran-server/model"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	log "github.com/sirupsen/logrus"
)

func NewBuild(ctx *gin.Context) {
	proj := ctx.Keys["proj"].(*model.Project)
	cfg := model.GetProjCfgByID(proj.ID)
	lk := model.ProjMutexMap.Lock(proj.ID)
	defer lk.Unlock()
	repo, err := git.PlainOpen(proj.Path)
	if err != nil {
		log.Errorf("git.PlainOpen(%s): %+v", proj.Path, err.Error())
		return
	}
	wt, _ := repo.Worktree()
	srcBr := "refs/heads/" + cfg.SrcBr
	trnBr := "refs/heads/" + cfg.TrnBr
	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(trnBr),
	})
	if err != nil {
		log.Errorf("wt.Checkout(%s): %+v", cfg.TrnBr, err.Error())
		return
	}
	files := model.ListProjFiles(proj.ID)
	wg := sync.WaitGroup{}
	wg.Add(len(files))
	for _, file := range files {
		go buildFile(&wg, &file, proj, cfg)
	}
	wg.Wait()
	wt.Add(".")
	hash, err := wt.Commit("new translations", &git.CommitOptions{
		All: false,
		Author: &object.Signature{
			Name:  config.APP.Name,
			Email: "noreply@" + config.APP.Domain,
			When:  time.Now(),
		},
	})
	if err != nil {
		log.Errorf("wt.Commit: %+v", err.Error())
		ctx.JSON(http.StatusOK, model.Response{
			Success: false,
			Code:    constant.ErrGitCommit,
		})
		return
	}
	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(srcBr),
	})
	if err != nil {
		log.Errorf("wt.Checkout(%s): %+v", cfg.SrcBr, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, model.Response{
		Success: true,
		Data: gin.H{
			"hash": hash.String(),
		}})
}

func buildFile(wg *sync.WaitGroup, file *model.ProjFile, proj *model.Project, cfg *model.ProjCfg) {
	defer wg.Done()
	sents := model.ListValidSentsOrderByOffset(file.ID)
	srcFile := file.ReadContent()
	for _, lang := range proj.TranslateLanguages {
		srcOff := 0
		buf := bytes.Buffer{}
		for _, sent := range sents {
			tran := sent.TopTran(lang.Code)
			content := []byte(sent.Content)
			if srcOff < sent.Offset {
				buf.Write(srcFile[srcOff:sent.Offset])
			}
			if tran != nil {
				content = []byte(tran.Content)
			}
			buf.Write(content)
			srcOff = sent.Offset + len(sent.Content)
		}
		if srcOff < len(srcFile) {
			buf.Write(srcFile[srcOff:])
		}
		filename := model.GenTrnFilesFromSrcFiles([]string{file.Path}, cfg.TrnReg, &lang, proj)[0]
		writeFile(path.Join(proj.Path, filename), buf.Bytes())
	}
}

func writeFile(filename string, data []byte) {
	dir := path.Dir(filename)
	os.MkdirAll(dir, os.ModePerm)
	os.WriteFile(filename, data, 0644)
}
