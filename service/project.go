package service

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"

	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/config"
	"github.com/gitran-com/gitran-server/constant"
	"github.com/gitran-com/gitran-server/model"
	"github.com/gitran-com/gitran-server/util"
	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
	log "github.com/sirupsen/logrus"
)

var (
	urlNameReg = "^[A-Za-z0-9-_]{1,32}$"
)

//validateURLName checks if a name is legal (to be in URL)
func validateURLName(name string) bool {
	ok, _ := regexp.Match(urlNameReg, []byte(name))
	return ok
}

func validateGitURL(url string) bool {
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

//GetOrgProjByName get an org project info
func GetOrgProjByName(ctx *gin.Context, owner string, name string) *model.Project {
	//TODO
	return nil
}

//GetProj get a project info
func GetProj(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	proj := model.GetProjByID(id)
	if proj == nil {
		ctx.JSON(http.StatusNotFound, util.Resp404)
		return
	}
	projInfo := model.GetProjInfoFromProj(proj)
	ctx.JSON(http.StatusOK, util.Response{
		Success: true,
		Data: gin.H{
			"proj_info": *projInfo,
		},
	})
}

//ListUserProj list all projects
func ListUserProj(ctx *gin.Context) {
	user_id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	ctx.JSON(http.StatusOK, util.Response{
		Success: true,
		Msg:     "",
		Data: gin.H{
			"proj_infos": model.GetProjInfosFromProjs(model.ListUserProj(user_id)),
		}})
}

//CreateUserProj create a new user project
func CreateUserProj(ctx *gin.Context) {
	name := ctx.PostForm("name")
	uri := ctx.PostForm("uri")
	desc := ctx.PostForm("desc")
	user := ctx.Keys["user"].(*model.User)
	gitURL := ctx.PostForm("git_url")
	src := ctx.PostForm("src_langs")
	trn := ctx.PostForm("trn_langs")
	srcLangs := model.ParseLangs(src)
	trnLangs := model.ParseLangs(trn)
	tp, _ := strconv.Atoi(ctx.PostForm("type"))
	var token *model.Token
	var token_id int64
	var err error
	if token_id, err = strconv.ParseInt(ctx.PostForm("token_id"), 10, 64); err == nil {
		token = model.GetTokenByID(token_id)
	} else {
		token = nil
	}
	if token == nil && tp == constant.TypeGithub {
		ctx.JSON(http.StatusBadRequest, util.Response{
			Success: false,
			Msg:     "Token不合法",
			Data:    nil,
			Code:    constant.ErrTokenIllegal,
		})
		return
	}
	if model.GetProjByURI(uri) != nil {
		ctx.JSON(http.StatusBadRequest, util.Response{
			Success: false,
			Msg:     "URI已被占用",
			Data:    nil,
			Code:    constant.ErrProjUriIllegal,
		})
		return
	}
	if !validateURLName(name) {
		ctx.JSON(http.StatusBadRequest, util.Response{
			Success: false,
			Msg:     "名字不合法",
			Data:    nil,
			Code:    constant.ErrProjNameIllegal,
		})
		return
	}
	if len(srcLangs) == 0 {
		ctx.JSON(http.StatusBadRequest, util.Response{
			Success: false,
			Msg:     "源语言不能为空",
			Data:    nil,
			Code:    constant.ErrProjSrcLangEmpty,
		})
		return
	}
	if !model.ValidateLangs(srcLangs) {
		ctx.JSON(http.StatusBadRequest, util.Response{
			Success: false,
			Msg:     "源语言不合法",
			Data:    nil,
			Code:    constant.ErrProjSrcLangIllegal,
		})
		return
	}
	if !model.ValidateLangs(trnLangs) {
		ctx.JSON(http.StatusBadRequest, util.Response{
			Success: false,
			Msg:     "目标语言不合法",
			Data:    nil,
			Code:    constant.ErrProjTrnLangIllegal,
		})
		return
	}
	if tp == constant.TypeGitURL && !validateGitURL(gitURL) {
		ctx.JSON(http.StatusBadRequest, util.Response{
			Success: false,
			Msg:     "Git URL不合法",
			Data:    nil,
			Code:    constant.ErrGitUrlIllegal,
		})
		return
	}
	fmt.Printf("token=%+v\n", token)
	if token != nil && token.OwnerID != user.ID {
		ctx.JSON(http.StatusBadRequest, util.Response{
			Success: false,
			Msg:     "Token不合法",
			Data:    nil,
			Code:    constant.ErrTokenIllegal,
		})
		return
	}
	proj, err := model.NewProj(&model.Project{
		Name:     name,
		OwnerID:  user.ID,
		TokenID:  token_id,
		Token:    token,
		Type:     tp,
		Status:   constant.ProjStatCreated,
		Desc:     desc,
		GitURL:   gitURL,
		Path:     config.ProjPath + "/" + uri + "/",
		SrcLangs: src,
		TrnLangs: trn,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.Response{
			Success: false,
			Msg:     err.Error(),
			Data:    nil,
		})
		return
	}
	go initProj(proj)
	ctx.JSON(http.StatusCreated, util.Response{
		Success: true,
		Msg:     "创建项目成功",
		Data:    nil,
	})
}

//initProj init a new project
func initProj(proj *model.Project) {
	if proj.Type == constant.TypeGitURL {
		_, err := git.PlainClone(proj.Path, false, &git.CloneOptions{
			URL:          proj.GitURL,
			Progress:     os.Stdout, //TODO: update progress in DB
			Depth:        1,
			SingleBranch: false,
		})
		if err == nil {
			model.UpdateProjStatus(proj, constant.ProjStatReady)
		} else {
			log.Warnf("git clone error : %v", err.Error())
		}
	} else if proj.Type == constant.TypeGithub {
		_, err := git.PlainClone(proj.Path, false, &git.CloneOptions{
			URL: proj.GitURL,
			Auth: &githttp.BasicAuth{
				Username: config.APP.Name,
				Password: proj.Token.AccessToken,
			},
			Progress:     os.Stdout,
			Depth:        1,
			SingleBranch: false,
		})
		if err == nil {
			model.UpdateProjStatus(proj, constant.ProjStatReady)
		} else {
			log.Warnf("git clone error : %v", err.Error())
		}
	} else {
		log.Errorf("initProj error: type %v has not been implemented", proj.Type)
		//TODO
	}
}

func getBrchFromRef(ref string) string {
	return strings.TrimPrefix(ref, "refs/heads/")
}

func ListProjBrch(ctx *gin.Context) {
	proj := ctx.Keys["project"].(*model.Project)
	repo, _ := git.PlainOpen(proj.Path)
	refs, _ := repo.Branches()
	var brchs []string
	refs.ForEach(func(r *plumbing.Reference) error {
		brchs = append(brchs, getBrchFromRef(string(r.Name())))
		// fmt.Printf("branch-name=%+v\n", getBrchFromRef(string(r.Name())))
		return nil
	})
	ctx.JSON(http.StatusOK, util.Response{
		Success: true,
		Data: gin.H{
			"branches": brchs,
		},
	})
}
