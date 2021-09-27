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
)

func pushGit(proj_id int64) {
	for {
		log.Infof("begin to push project %v", proj_id)
		lk := model.ProjMutexMap.Lock(proj_id)
		defer lk.Unlock()
		proj := model.GetProjByID(proj_id)
		cfg := model.GetProjCfgByID(proj_id)
		if proj == nil {
			log.Warnf("project %v not found when pushing", proj_id)
			return
		}
		if cfg == nil {
			log.Warnf("projcfg %v not found when pushing", proj_id)
			return
		}
		if cfg.PushGap == 0 {
			log.Warnf("project %v no longer needs push", proj_id)
			return
		}
		if cfg.PushStatus == constant.SyncStatDoing {
			log.Warnf("project %v push conflicted", proj_id)
			return
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
			log.Warnf("prject %v checkout failed", proj_id)
			model.UpdateProjCfgPushStatus(cfg, constant.SyncStatFail)
			return
		}
		tk := proj.Token
		if tk == "" {
			log.Warnf("project %v no token", proj_id)
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
			log.Warnf("project %v push failed : %v", proj_id, err.Error())
			model.UpdateProjCfgPushStatus(cfg, constant.SyncStatFail)
			return
		}
		log.Infof("project %v push successfully", proj_id)
		model.UpdateProjCfgPushStatus(cfg, constant.SyncStatSucc)
		time.Sleep(time.Minute * time.Duration(cfg.PushGap))
	}
}

func pullGit(proj_id int64) {
	for {
		log.Infof("begin to pull project %v", proj_id)
		lk := model.ProjMutexMap.Lock(proj_id)
		defer lk.Unlock()
		proj := model.GetProjByID(proj_id)
		cfg := model.GetProjCfgByID(proj_id)
		if proj == nil {
			log.Warnf("project %v not found when pulling", proj_id)
			return
		}
		if cfg == nil {
			log.Warnf("projcfg %v not found when pulling", proj_id)
			return
		}
		if cfg.PullGap == 0 {
			log.Warnf("project %v no longer needs pull", proj_id)
			return
		}
		if cfg.PullStatus == constant.SyncStatDoing {
			log.Warnf("project %v pull conflicted", proj_id)
			return
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
			log.Warnf("project %v checkout failed", proj_id)
			model.UpdateProjCfgPullStatus(cfg, constant.SyncStatFail)
			return
		}
		tk := proj.Token
		if proj.Type == constant.ProjTypeGithub && tk == "" {
			log.Warnf("project %v no token", proj_id)
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
			log.Warnf("project %v pull failed : %v", proj_id, err.Error())
			model.UpdateProjCfgPullStatus(cfg, constant.SyncStatFail)
			return
		}
		log.Infof("project %v pull successfully", proj_id)
		model.UpdateProjCfgPullStatus(cfg, constant.SyncStatSucc)
		time.Sleep(time.Minute * time.Duration(cfg.PullGap))
	}
}
