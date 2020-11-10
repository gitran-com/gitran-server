package service

import (
	log "github.com/sirupsen/logrus"
	"github.com/wzru/gitran-server/model"
)

//Init init the service
func Init() error {
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
	return nil
}
