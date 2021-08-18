package service

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/gitran-com/gitran-server/config"
	"github.com/gitran-com/gitran-server/model"
	"github.com/gorilla/websocket"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	log "github.com/sirupsen/logrus"
)

var (
	GithubUserProvider *github.Provider
	GithubRepoProvider *github.Provider
	GithubRoute        string
	upGrader           websocket.Upgrader
)

//Init init the service
func Init() error {
	pullSchd.StartAsync()
	pushSchd.StartAsync()
	cfgs := model.ListSyncProjCfg()
	for _, cfg := range cfgs {
		proj := model.GetProjByID(cfg.ID)
		if proj == nil {
			log.Warnf("proj %v not found", cfg.ID)
		} else {
			if cfg.PullGap != 0 {
				// log.Infof("add pull task %+v", cfg)
				_, err := pullSchd.Every(uint64(cfg.PullGap)).Minutes().Do(pullGit, &cfg)
				if err != nil {
					log.Warnf("add pull task err : %+v", err.Error())
				}
			}
			if cfg.PushGap != 0 {
				// log.Infof("add push task %+v", cfg)
				_, err := pushSchd.Every(uint64(cfg.PushGap)).Minutes().Do(pushGit, &cfg)
				if err != nil {
					log.Warnf("add push task err : %+v", err.Error())
				}
			}
		}
	}
	InitGithubAuth()
	InitWebsocket()
	go InitProj()
	return nil
}

func InitProj() {
	projs := model.ListUninitProj()
	wg := sync.WaitGroup{}
	wg.Add(len(projs))
	for _, proj := range projs {
		go func(wg *sync.WaitGroup, p *model.Project) {
			defer wg.Done()
			p.Init()
		}(&wg, &proj)
	}
	wg.Wait()
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

func InitWebsocket() {
	if config.IsDebug {
		upGrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return false
			},
		}
	} else {
		upGrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
	}
}
