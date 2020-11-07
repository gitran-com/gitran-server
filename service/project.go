package service

import (
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
	log "github.com/sirupsen/logrus"
	"github.com/wzru/gitran-server/config"
	"github.com/wzru/gitran-server/constant"
	"github.com/wzru/gitran-server/middleware"
	"github.com/wzru/gitran-server/model"
)

var (
	urlNameReg = "^[A-Za-z0-9-_]{1,32}$"
)

//checkURLName checks if a name is legal (to be in URL)
func checkURLName(name string) bool {
	ok, _ := regexp.Match(urlNameReg, []byte(name))
	return ok
}

func checkGitURL(url string) bool {
	rmt := git.NewRemote(memory.NewStorage(), &gitconfig.RemoteConfig{
		URLs: []string{url},
	})
	if _, err := rmt.List(&git.ListOptions{}); err != nil {
		log.Warnf("Git url error : %v", err.Error())
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
	owner := ctx.Param("owner")
	name := ctx.Param("project")
	proj := GetUserProjByName(ctx, owner, name)
	if proj == nil {
		proj = GetOrgProjByName(ctx, owner, name)
	}
	if proj == nil {
		ctx.JSON(http.StatusNotFound, model.Result404)
		return
	}
	projInfo := model.GetProjInfoFromProj(proj)
	ctx.JSON(http.StatusOK, model.Result{
		Success: false,
		Msg:     "",
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
				"projects:": model.GetProjInfosFromProjs(model.ListProjFromUser(usr, middleware.HasUserPermission(ctx, usr.ID))),
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
	name := ctx.PostForm("name")
	ot := constant.OwnerUsr
	desc := ctx.PostForm("desc")
	userID, _ := strconv.Atoi(ctx.GetString("user-id"))
	userName := ctx.GetString("user-name")
	isPrvt := ctx.PostForm("is_private") == "true"
	gitURL := ctx.PostForm("git_url")
	src := ctx.PostForm("src_langs")
	trn := ctx.PostForm("trn_langs")
	srcLangs := model.GetLangsFromString(src)
	trnLangs := model.GetLangsFromString(trn)
	tp, _ := strconv.Atoi(ctx.PostForm("type"))
	if !checkURLName(name) {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     "名字不合法",
			Data:    nil,
		})
		return
	}
	if len(srcLangs) == 0 {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     "源语言不能为空",
			Data:    nil,
		})
		return
	}
	if !model.CheckLangs(srcLangs) {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     "源语言不合法",
			Data:    nil,
		})
		return
	}
	if !model.CheckLangs(trnLangs) {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     "目标语言不合法",
			Data:    nil,
		})
		return
	}
	if !checkGitURL(gitURL) {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     "Git URL不合法",
			Data:    nil,
		})
		return
	}
	proj := &model.Project{
		Name:      name,
		Desc:      desc,
		OwnerID:   uint64(userID),
		OwnerType: uint8(ot),
		IsPrivate: isPrvt,
		Type:      uint8(tp),
		Status:    uint8(constant.ProjStatCreated),
		GitURL:    gitURL,
		Path:      config.ProjPath + userName + "/" + name + "/",
		SrcLangs:  src,
		TrnLangs:  trn,
	}
	if findProj := model.GetProjByOwnerIDName(proj.OwnerID, proj.Name, true); findProj != nil {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     "项目已存在",
			Data:    nil,
		})
		return
	}
	proj, err := model.NewProj(proj)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     err.Error(),
			Data:    nil,
		})
		return
	}
	ctx.JSON(http.StatusCreated, model.Result{
		Success: true,
		Msg:     "创建项目成功",
		Data:    nil,
	})
	go initProj(proj)
}

//initUserProj init a new user project
func initUserProj(proj *model.Project) {
	if proj.Type == constant.ProjGithub {
		_, err := git.PlainClone(proj.Path, false, &git.CloneOptions{
			URL:          proj.GitURL,
			Progress:     os.Stdout,
			Depth:        1,
			SingleBranch: false,
		})
		if err == nil {
			model.UpdateProjStatus(proj, constant.ProjStatInit)
		} else {
			log.Warnf("git clone error : %v", err.Error())
		}
	} else {
		//TODO
	}
}

//initOrgProj init a new org project
func initOrgProj(proj *model.Project) {
	//TODO
}

//initProj init a new project
func initProj(proj *model.Project) {
	if proj.OwnerType == constant.OwnerUsr {
		initUserProj(proj)
	} else if proj.OwnerType == constant.OwnerOrg {
		initOrgProj(proj)
	}
}

//CreateOrgProj create a new organization project
func CreateOrgProj(ctx *gin.Context) {
	//TODO
}
