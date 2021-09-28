package service

import (
	"os"
	"os/exec"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gitran-com/gitran-server/constant"
	"github.com/gitran-com/gitran-server/model"
)

func doPush(proj_id int64) (gap int, ok bool) {
	lk := model.ProjMutexMap.Lock(proj_id)
	log.Infof("begin to push project %v", proj_id)
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
	gap = int(cfg.PushGap)
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
	ck := exec.Command("git", "checkout", cfg.TrnBr)
	ck.Dir = proj.Path
	if err := ck.Run(); err != nil {
		log.Warnf("project %v checkout failed: %v", proj_id, err.Error())
		model.UpdateProjCfgPushStatus(cfg, constant.SyncStatFail)
		return
	}
	push := exec.Command("git", "push", "origin", cfg.TrnBr)
	push.Dir = proj.Path
	push.Stderr = os.Stderr
	if err := push.Run(); err != nil {
		log.Warnf("project %v push failed: %v; will retry...", proj_id, err.Error())
		model.UpdateProjCfgPushStatus(cfg, constant.SyncStatFail)
		ok = true
		return
	}
	log.Infof("project %v push successfully", proj_id)
	model.UpdateProjCfgPushStatus(cfg, constant.SyncStatSucc)
	ok = true
	return
}

func doPull(proj_id int64) (gap int, ok bool) {
	lk := model.ProjMutexMap.Lock(proj_id)
	log.Infof("begin to pull project %v", proj_id)
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
	gap = int(cfg.PullGap)
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
	ck := exec.Command("git", "checkout", cfg.SrcBr)
	ck.Dir = proj.Path
	if err := ck.Run(); err != nil {
		log.Warnf("project %v checkout failed", proj_id)
		model.UpdateProjCfgPushStatus(cfg, constant.SyncStatFail)
		return
	}
	pull := exec.Command("git", "pull")
	pull.Dir = proj.Path
	pull.Stderr = os.Stderr
	if err := pull.Run(); err != nil {
		log.Warnf("project %v pull failed : %v; will retry...", proj_id, err.Error())
		model.UpdateProjCfgPullStatus(cfg, constant.SyncStatFail)
		ok = true
		return
	}
	log.Infof("project %v pull successfully", proj_id)
	model.UpdateProjCfgPullStatus(cfg, constant.SyncStatSucc)
	ok = true
	return
}

func pushGit(proj_id int64) {
	var (
		ok  bool
		gap int
	)
	for {
		if gap, ok = doPush(proj_id); !ok {
			return
		}
		time.Sleep(time.Minute * time.Duration(gap))
	}
}

func pullGit(proj_id int64) {
	var (
		ok  bool
		gap int
	)
	for {
		if gap, ok = doPull(proj_id); !ok {
			return
		}
		time.Sleep(time.Minute * time.Duration(gap))
	}
}
