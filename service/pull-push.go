package service

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	log "github.com/sirupsen/logrus"

	"github.com/go-co-op/gocron"
	"github.com/wzru/gitran-server/constant"
	"github.com/wzru/gitran-server/model"
)

var (
	pullSchd = gocron.NewScheduler(time.UTC)
	pushSchd = gocron.NewScheduler(time.UTC)
)

func pushGit(cfg *model.ProjCfg) {
	fmt.Printf("begin to push...")
	beg := time.Now().Unix()
	projMutexMap.Lock(cfg.ProjID)
	defer projMutexMap.Unlock(cfg.ProjID)
	proj := model.GetProjByID(cfg.ProjID)
	if proj == nil {
		return
	}
	cfg = model.GetProjCfgByID(cfg.ID)
	if cfg == nil {
		return
	}
	if cfg.LastPushAt.Unix() >= beg {
		log.Warnf("project %v, config %v PUSH time out", proj.ID, cfg.ID)
		return
	}
	if cfg.PushStatus == constant.SyncStatDoing {
		log.Warnf("project %v, config %v last PUSH aborted", proj.ID, cfg.ID)
	} else {
		model.UpdateProjCfgPushStatus(cfg, constant.SyncStatDoing)
	}
	brName := "refs/heads/" + cfg.TrnBr
	repo, _ := git.PlainOpen(proj.Path)
	wt, _ := repo.Worktree()
	err := wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(brName),
	})
	if err != nil {
		log.Warnf("%v checkout failed", proj.Path)
		model.UpdateProjCfgPushStatus(cfg, constant.SyncStatFail)
		return
	}
	err = repo.Push(&git.PushOptions{RemoteName: "origin"})
	if err != nil {
		log.Warnf("%v push failed", proj.Path)
		model.UpdateProjCfgPushStatus(cfg, constant.SyncStatFail)
		return
	}
	log.Infof("%v push successfully", proj.Path)
	model.UpdateProjCfgPushStatus(cfg, constant.SyncStatSucc)
}

func pullGit(cfg *model.ProjCfg) {
	fmt.Printf("begin to pull...")
	beg := time.Now().Unix()
	projMutexMap.Lock(cfg.ProjID)
	defer projMutexMap.Unlock(cfg.ProjID)
	proj := model.GetProjByID(cfg.ProjID)
	if proj == nil {
		return
	}
	cfg = model.GetProjCfgByID(cfg.ID)
	if cfg == nil {
		return
	}
	if cfg.LastPullAt.Unix() >= beg {
		log.Warnf("project %v, config %v PULL time out", proj.ID, cfg.ID)
		return
	}
	if cfg.PullStatus == constant.SyncStatDoing {
		log.Warnf("project %v, config %v last PULL aborted", proj.ID, cfg.ID)
	} else {
		model.UpdateProjCfgPullStatus(cfg, constant.SyncStatDoing)
	}
	brName := "refs/heads/" + cfg.SrcBr
	repo, _ := git.PlainOpen(proj.Path)
	wt, _ := repo.Worktree()
	err := wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(brName),
	})
	if err != nil {
		log.Warnf("%v checkout failed", proj.Path)
		model.UpdateProjCfgPullStatus(cfg, constant.SyncStatFail)
		return
	}
	err = wt.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil {
		log.Warnf("%v pull failed", proj.Path)
		model.UpdateProjCfgPullStatus(cfg, constant.SyncStatFail)
		return
	}
	log.Infof("%v pull successfully", proj.Path)
	model.UpdateProjCfgPullStatus(cfg, constant.SyncStatSucc)
}
