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
	initGithubAuth()
	initWebsocket()
	go initSync()
	go initProj()
	return nil
}

func initProj() {
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

func initGithubAuth() {
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

func initWebsocket() {
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

func initSync() {
	cfgs := model.ListSyncProjCfg()
	for _, cfg := range cfgs {
		if cfg.PushGap != 0 {
			go pushGit(cfg.ID)
		}
		if cfg.PullGap != 0 {
			go pullGit(cfg.ID)
		}
	}
}
