package service

import (
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
	log "github.com/sirupsen/logrus"
	"github.com/wzru/gitran-server/config"
	"github.com/wzru/gitran-server/constant"
	"github.com/wzru/gitran-server/middleware"
	"github.com/wzru/gitran-server/model"
	"github.com/wzru/gitran-server/util"
)

var (
	nameReg = `^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*$`
)

//checkURLName checks if a name is legal (to be in URL)
func checkURLName(name string) bool {
	if len(name) > 32 {
		return false
	}
	ok, _ := regexp.Match(nameReg, []byte(name))
	return ok
}

func checkGitURL(url string) bool {
	rmt := git.NewRemote(memory.NewStorage(), &gitconfig.RemoteConfig{
		URLs: []string{url},
	})
	if _, err := rmt.List(&git.ListOptions{}); err != nil {
		log.Warnf("git url '%s' error : %v", url, err.Error())
		return false
	}
	return true
}

func createGitProj(ctx *gin.Context) error {
	return nil
}

//GetUserProjByName get a user project info
func GetUserProjByName(ctx *gin.Context, owner string, name string) *model.Project {
	user := model.GetUserByName(owner)
	if user == nil {
		return nil
	}
	self := (middleware.HasUserPermission(ctx, user.ID))
	return model.GetProjByOwnerIDName(user.ID, name, self)
}

//GetOrgProjByName get an org project info
func GetOrgProjByName(ctx *gin.Context, owner string, name string) *model.Project {
	//TODO
	return nil
}

//GetProj get a project info
func GetProj(ctx *gin.Context) {
	user, _, tp := model.GetOwnerByName(ctx.Param("owner"))
	name := ctx.Param("project")
	var proj *model.Project
	if tp == constant.OwnerUsr {
		proj = model.GetProjByOwnerIDName(user.ID, name, middleware.HasUserPermission(ctx, user.ID))
	} else if tp == constant.OwnerOrg {
		ctx.JSON(http.StatusNotFound, model.Result404)
		return
	}
	if proj == nil {
		ctx.JSON(http.StatusNotFound, model.Result404)
		return
	}
	projInfo := model.GetProjInfoFromProj(proj)
	ctx.JSON(http.StatusOK, model.Result{
		Success: true,
		Data: gin.H{
			"proj_info": *projInfo,
		},
	})
}

//ListProj list all projects
func ListProj(ctx *gin.Context) {
	usr, org, tp := model.GetOwnerByName(ctx.Param("owner"))
	if tp == constant.OwnerUsr {
		if usr == nil {
			ctx.JSON(http.StatusNotFound, model.Result404)
			return
		}
		ctx.JSON(http.StatusOK, model.Result{
			Success: true,
			Msg:     "",
			Data: gin.H{
				"projects:": model.GetProjInfosFromProjs(model.ListProjByUserID(usr.ID, middleware.HasUserPermission(ctx, usr.ID))),
			}})
		return
	} else {
		//TODO
		if org == nil {
			ctx.JSON(http.StatusNotFound, model.Result404)
			return
		}
	}
}

//CreateUserProj create a new user project
func CreateUserProj(ctx *gin.Context) {
	user := ctx.Keys["user"].(*model.User)
	name := ctx.PostForm("name")
	desc := ctx.PostForm("desc")
	prv := ctx.PostForm("private") == "true"
	importURL := ctx.PostForm("import_url")
	repoID := util.ParseUint64(ctx.PostForm("repo_id"))
	src := ctx.PostForm("src_langs")
	trn := ctx.PostForm("trn_langs")
	srcLangs := model.GetLangsFromString(src)
	trnLangs := model.GetLangsFromString(trn)
	tp := util.ParseUint8(ctx.PostForm("type"))
	if !checkURLName(name) {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     "名字不合法",
			Data:    nil,
			Code:    constant.ErrorProjNameIllegal,
		})
		return
	}
	if len(srcLangs) == 0 {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     "源语言不能为空",
			Data:    nil,
			Code:    constant.ErrorProjSrcLangEmpty,
		})
		return
	}
	if !model.CheckLangs(srcLangs) {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     "源语言不合法",
			Data:    nil,
			Code:    constant.ErrorProjNameIllegal,
		})
		return
	}
	if !model.CheckLangs(trnLangs) {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     "目标语言不合法",
			Data:    nil,
			Code:    constant.ErrorProjTrnLangIllegal,
		})
		return
	}
	url, err := util.ParseGitURL(importURL)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     "Git URL不合法",
			Data:    nil,
			Code:    constant.ErrorProjGitURLIllegal,
		})
		return
	}
	importURL = util.GenHTTPSGitURL(url.Host, url.Path)
	if tp == constant.TypeGitURL {
		if !checkGitURL(importURL) {
			ctx.JSON(http.StatusBadRequest, model.Result{
				Success: false,
				Msg:     "Git URL不合法",
				Data:    nil,
				Code:    constant.ErrorProjGitURLIllegal,
			})
			return
		}
	} else if tp == constant.TypeGithub {
		tk := model.GetTokenByOwnerID(user.ID, constant.TypeGithub, constant.TokenRepo)
		if tk == nil {
			ctx.JSON(http.StatusBadRequest, model.Result{
				Success: false,
				Msg:     "GitHub unauthorized",
				Data:    nil,
				Code:    constant.ErrorGithubUnauthorized,
			})
			return
		}
		repo := getGithubRepoByTokenID(repoID, tk)
		if repo == nil {
			ctx.JSON(http.StatusBadRequest, model.Result404)
			return
		}
		importURL = repo.CloneURL
	}
	proj := &model.Project{
		Name:      name,
		Desc:      desc,
		OwnerID:   user.ID,
		OwnerType: constant.OwnerUsr,
		Private:   prv,
		Type:      tp,
		Status:    constant.ProjStatCreated,
		GitURL:    importURL,
		RepoID:    repoID,
		Path:      config.ProjPath + user.Login + "/" + name + "/",
		SrcLangs:  src,
		TrnLangs:  trn,
	}
	if findProj := model.GetProjByOwnerIDName(proj.OwnerID, proj.Name, true); findProj != nil {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     "项目已存在",
			Data:    nil,
			Code:    constant.ErrorProjNameExists,
		})
		return
	}
	proj, err = model.NewProj(proj)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.Result{
			Success: false,
			Msg:     err.Error(),
			Data:    nil,
		})
		return
	}
	go initUserProj(proj)
	ctx.JSON(http.StatusCreated, model.Result{
		Success: true,
		Msg:     "创建项目成功",
		Data:    nil,
	})
}

//initUserProj init a new user project
func initUserProj(proj *model.Project) {
	if proj.Type == constant.TypeGithub {
		if _, err := os.Stat(proj.Path); !os.IsNotExist(err) {
			log.Warnf("remove path %v", proj.Path)
			os.RemoveAll(proj.Path)
		}
		_, err := git.PlainClone(proj.Path, false, &git.CloneOptions{
			URL:      proj.GitURL,
			Progress: os.Stdout,
			Depth:    1,
		})
		if err == nil {
			model.UpdateProjStatus(proj, constant.ProjStatInit)
		} else {
			log.Warnf("git clone %+v into %+v ERROR : %+v", proj.GitURL, proj.Path, err.Error())
		}
	} else if proj.Type == constant.TypeGitURL {
		if _, err := os.Stat(proj.Path); !os.IsNotExist(err) {
			log.Warnf("remove path %v", proj.Path)
			os.RemoveAll(proj.Path)
		}
		_, err := git.PlainClone(proj.Path, false, &git.CloneOptions{
			URL:      proj.GitURL,
			Progress: os.Stdout,
		})
		if err == nil {
			model.UpdateProjStatus(proj, constant.ProjStatInit)
		} else {
			log.Warnf("git clone %+v into %+v ERROR : %+v", proj.GitURL, proj.Path, err.Error())
		}
	} else {
		model.UpdateProjStatus(proj, constant.ProjStatInit)
	}
}

//initOrgProj init a new org project
func initOrgProj(proj *model.Project) {
	//TODO
}

//CreateOrgProj create a new organization project
func CreateOrgProj(ctx *gin.Context) {
	//TODO
}

//ListUserPubProj list public projects from a user
func ListUserPubProj(ctx *gin.Context) {
	user := ctx.Keys["user"].(*model.User)
	ctx.JSON(http.StatusOK, model.Result{
		Success: true,
		Data: gin.H{
			"proj_infos": model.GetProjInfosFromProjs(model.ListProjByUserID(user.ID, false)),
		}})
}

//ListAuthUserProj list public projects from a auth user
func ListAuthUserProj(ctx *gin.Context) {
	user := ctx.Keys["user"].(*model.User)
	ctx.JSON(http.StatusOK, model.Result{
		Success: true,
		Data: gin.H{
			"proj_infos": model.GetProjInfosFromProjs(model.ListProjByUserID(user.ID, true)),
		}})
}

//L
func ListUserProjBrch(ctx *gin.Context) {

}
