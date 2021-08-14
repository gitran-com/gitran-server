package service

import (
	"fmt"
	"net/url"

	"github.com/gitran-com/gitran-server/config"
	"github.com/gitran-com/gitran-server/model"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	log "github.com/sirupsen/logrus"
)

var (
	GithubUserProvider *github.Provider
	GithubRepoProvider *github.Provider
	GithubRoute        string
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
				_, err := pullSchd.Every(uint64(cfg.PullItv)).Minutes().Do(pullGit, &cfg)
				if err != nil {
					log.Warnf("add pull task err : %+v", err.Error())
				}
			}
			if cfg.PushItv != 0 {
				// log.Infof("add push task %+v", cfg)
				_, err := pushSchd.Every(uint64(cfg.PushItv)).Minutes().Do(pushGit, &cfg)
				if err != nil {
					log.Warnf("add push task err : %+v", err.Error())
				}
			}
		}
	}
	InitProj()
	InitGithubAuth()
	return nil
}

func InitProj() {
	projs := model.ListUninitProj()
	for _, proj := range projs {
		go proj.Init()
	}
}

func InitGithubAuth() {
	if config.Github.Enable {
		url, err := url.Parse(config.APP.URL)
		if err != nil {
			log.Errorf("service/init parse url error: %+v", err.Error())
		}
		if url.Scheme != "http" && url.Scheme != "https" {
			url.Scheme = "http"
		}
		fmt.Printf("url[%s]=%+v\n", config.APP.URL, *url)
		GithubRoute = fmt.Sprintf("%s://%s%s/api/v1/auth/github/", url.Scheme, url.Host, config.APP.APIPrefix)
		GithubUserProvider = github.New(config.Github.ClientID, config.Github.ClientSecret, GithubRoute+"login", "user")
		GithubUserProvider.SetName("github-user")
		GithubRepoProvider = github.New(config.Github.ClientID, config.Github.ClientSecret, GithubRoute+"import", "repo")
		GithubRepoProvider.SetName("github-repo")
		goth.UseProviders(GithubUserProvider, GithubRepoProvider)
	}
}
