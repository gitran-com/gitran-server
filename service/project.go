package service

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
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

func checkLang(lang *config.Lang) bool {
	for _, cfgLang := range config.Langs {
		if cfgLang.Code == lang.Code {
			return true
		}
	}
	return false
}

func checkLangs(langs []config.Lang) bool {
	for _, lang := range langs {
		ok := checkLang(&lang)
		if !ok {
			return false
		}
	}
	return true
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
	user := model.GetUserByUname(owner)
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
	// model.GetUserByUname("")
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
	// syncTime, _ := strconv.Atoi(ctx.PostForm("sync_time"))
	src := ctx.PostForm("src_langs")
	tgt := ctx.PostForm("tgt_langs")
	srcLangs := model.GetLangsFromString(src)
	tgtLangs := model.GetLangsFromString(tgt)
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
	if !checkLangs(srcLangs) {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     "源语言不合法",
			Data:    nil,
		})
		return
	}
	if !checkLangs(tgtLangs) {
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
		// SyncTime:  uint64(syncTime),
		SrcLangs: src,
		TgtLangs: tgt,
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

//GetUserProjCfg get a user project config
func GetUserProjCfg(ctx *gin.Context) {
	//TODO
}

//CreateUserProjCfg create a new project config
func CreateUserProjCfg(ctx *gin.Context) {
	// projID, _ := strconv.Atoi(ctx.GetString("proj-id"))
	proj := ctx.Keys["project"].(*model.Project)
	srcBrName := "refs/heads/" + ctx.PostForm("src_branch")
	tgtBrName := "refs/heads/" + ctx.PostForm("tgt_branch")
	syncTime, _ := strconv.ParseUint(ctx.PostForm("sync_time"), 10, 64)
	pushTrans := ctx.PostForm("push_trans") == "true"
	filename := ctx.PostForm("filename")
	// fmt.Printf("srcbr=%+v\n", srcBrName)
	repo, err := git.PlainOpen(proj.Path)
	wt, _ := repo.Worktree()
	//先切换到src分支
	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(srcBrName),
	})
	if err != nil {
		fmt.Println("HERE1")
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     err.Error(),
			Data:    nil,
		})
		return
	}
	//然后新建tgt分支
	srcHead, _ := repo.Head()
	ref := plumbing.NewHashReference(plumbing.ReferenceName(tgtBrName), srcHead.Hash())
	if err := repo.Storer.SetReference(ref); err != nil {
		fmt.Println("HERE2")
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     err.Error(),
			Data:    nil,
		})
		return
	}
	projCfg := &model.ProjCfg{
		ProjID:    proj.ID,
		SrcBr:     ctx.PostForm("src_branch"),
		TgtBr:     ctx.PostForm("tgt_branch"),
		SyncTime:  syncTime,
		PushTrans: pushTrans,
		FileName:  filename,
	}
	projCfg, err = model.NewProjCfg(projCfg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     err.Error(),
			Data:    nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, model.Result{
		Success: true,
		Msg:     "项目分支配置新建成功",
		Data:    nil,
	})
}
