package service

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wzru/gitran-server/config"
	"github.com/wzru/gitran-server/model"
)

var (
	urlNameReg = "^[A-Za-z0-9-]{1,32}$"
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

func createGitProj(ctx *gin.Context) error {
	return nil
}

//GetProj get project info
func GetProj(ctx *gin.Context) {

}

//CreateProj creates new project
func CreateProj(ctx *gin.Context) {
	name := ctx.PostForm("name")
	isUsers := ctx.PostForm("is_users") == "true"
	desc := ctx.PostForm("desc")
	orgID, _ := strconv.Atoi(ctx.PostForm("org_id"))
	userID, _ := strconv.Atoi(ctx.GetString("user-id"))
	isPrvt := ctx.PostForm("is_private") == "true"
	isGit := ctx.PostForm("is_git") == "true"
	gitURL := ctx.PostForm("git_url")
	gitBranch := ctx.PostForm("git_branch")
	syncTime, _ := strconv.Atoi(ctx.PostForm("sync_time"))
	src := ctx.PostForm("src_langs")
	tgt := ctx.PostForm("tgt_langs")
	srcLangs := model.GetLangsFromString(src)
	tgtLangs := model.GetLangsFromString(tgt)
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
	proj := &model.Project{
		Name:      name,
		Desc:      desc,
		IsUsers:   isUsers,
		IsPrivate: isPrvt,
		IsGit:     isGit,
		GitURL:    gitURL,
		GitBranch: gitBranch,
		SyncTime:  uint64(syncTime),
		SrcLangs:  src,
		TgtLangs:  tgt,
	}
	if isUsers {
		proj.OwnerID = uint64(userID)
	} else {
		//TODO: 判断是否属于org
		proj.OwnerID = uint64(orgID)
	}
	if findProj := model.GetProjByOwnerName(proj.OwnerID, proj.Name); findProj != nil {
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
	} else {
		ctx.JSON(http.StatusOK, model.Result{
			Success: true,
			Msg:     "创建项目成功",
			Data:    nil,
		})
	}
}
