package service

import (
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	log "github.com/sirupsen/logrus"
	"github.com/wzru/gitran-server/config"
	"github.com/wzru/gitran-server/model"
)

func initSync() error {
	pullSchd.StartAsync()
	pushSchd.StartAsync()
	cfgs := model.ListSyncProjCfg()
	for _, cfg := range cfgs {
		proj := model.GetProjByID(cfg.ProjID)
		if proj == nil {
			log.Warnf("proj %v not found", cfg.ProjID)
		} else {
			if cfg.PullItv != 0 {
				// log.Infof("add pull task %+v", cfg)
				_, err := pullSchd.Every(uint64(cfg.PullItv)).Minutes().SetTag([]string{proj.Path}).Do(pullGit, &cfg)
				if err != nil {
					log.Warnf("add pull task err : %+v", err.Error())
				}
			}
			if cfg.PushItv != 0 {
				// log.Infof("add push task %+v", cfg)
				_, err := pushSchd.Every(uint64(cfg.PushItv)).Minutes().SetTag([]string{proj.Path}).Do(pushGit, &cfg)
				if err != nil {
					log.Warnf("add push task err : %+v", err.Error())
				}
			}
		}
	}
	if config.Github.Enable {
		goth.UseProviders(github.New(config.Github.ClientID, config.Github.ClientSecret, config.Github.CallbackURL, ""))
	}
	return nil
}

//初始化未初始化的项目
func initProjInit() error {
	projs := model.ListProjByStatus(constant.ProjStatCreated)
	for i := range projs {
		go initUserProj(&projs[i])
	}
	return nil
}

//Init init the service
func Init() error {
	if err := initSync(); err != nil {
		return err
	}
	if err := initProjInit(); err != nil {
		return err
	}
	return nil
}
