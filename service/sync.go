package service

import (
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	log "github.com/sirupsen/logrus"

	"github.com/gitran-com/gitran-server/config"
	"github.com/gitran-com/gitran-server/constant"
	"github.com/gitran-com/gitran-server/model"
	"github.com/go-co-op/gocron"
)

var (
	pullSchd = gocron.NewScheduler(time.UTC)
	pushSchd = gocron.NewScheduler(time.UTC)
)

func pushGit(cfg *model.ProjCfg) {
	// log.Infof("begin to push project %v", cfg.ID)
	beg := time.Now().Unix()
	lk := projMutexMap.Lock(cfg.ID)
	defer lk.Unlock()
	proj := model.GetProjByID(cfg.ID)
	if proj == nil {
		log.Warnf("project %v not found when pushing", proj.ID)
		return
	}
	cid := cfg.ID
	cfg = model.GetProjCfgByID(cid)
	if cfg == nil {
		log.Warnf("project config %v not found when pushing", cfg.ID)
		return
	}
	if cfg.LastPushAt.Unix() >= beg {
		log.Warnf("project config not found when pushing", cid)
		return
	}
	if cfg.PushStatus == constant.SyncStatDoing {
		log.Warnf("project %v, config %v last push aborted", proj.ID, cfg.ID)
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
	tk := proj.Token
	if tk == "" {
		log.Warnf("%v get token failed", proj.Path)
		model.UpdateProjCfgPushStatus(cfg, constant.SyncStatFail)
		return
	}
	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: config.APP.Name,
			Password: tk,
		}})
	if err != nil && err.Error() != constant.ErrGitUpToDate {
		log.Warnf("%v push failed : %v", proj.Path, err.Error())
		model.UpdateProjCfgPushStatus(cfg, constant.SyncStatFail)
		return
	}
	log.Infof("%v push successfully", proj.Path)
	model.UpdateProjCfgPushStatus(cfg, constant.SyncStatSucc)
}

func pullGit(cfg *model.ProjCfg) {
	// log.Infof("begin to pull project %v", cfg.ID)
	beg := time.Now().Unix()
	lk := projMutexMap.Lock(cfg.ID)
	defer lk.Unlock()
	proj := model.GetProjByID(cfg.ID)
	if proj == nil {
		log.Warnf("project %v not found when pulling", proj.ID)
		return
	}
	cid := cfg.ID
	cfg = model.GetProjCfgByID(cid)
	if cfg == nil {
		log.Warnf("project config not found when pulling", cid)
		return
	}
	if cfg.LastPullAt.Unix() >= beg {
		log.Warnf("project %v, config %v pull time out", proj.ID, cfg.ID)
		return
	}
	if cfg.PullStatus == constant.SyncStatDoing {
		log.Warnf("project %v, config %v last pull aborted", proj.ID, cfg.ID)
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
	tk := proj.Token
	if proj.Type == constant.ProjTypeGithub && tk == "" {
		log.Warnf("%v get token failed", proj.Path)
		model.UpdateProjCfgPushStatus(cfg, constant.SyncStatFail)
		return
	}
	if tk != "" {
		err = wt.Pull(&git.PullOptions{
			RemoteName:   "origin",
			SingleBranch: true,
			Auth: &http.BasicAuth{
				Username: config.APP.Name,
				Password: tk,
			}})
	} else {
		err = wt.Pull(&git.PullOptions{
			RemoteName:   "origin",
			SingleBranch: true,
		})
	}
	if err != nil && err.Error() != constant.ErrGitUpToDate {
		log.Warnf("%v pull failed : %v", proj.Path, err.Error())
		model.UpdateProjCfgPullStatus(cfg, constant.SyncStatFail)
		return
	}
	log.Infof("%v pull successfully", proj.Path)
	model.UpdateProjCfgPullStatus(cfg, constant.SyncStatSucc)
}
