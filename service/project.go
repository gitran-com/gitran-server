package service

import (
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"

	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/config"
	"github.com/gitran-com/gitran-server/constant"
	"github.com/gitran-com/gitran-server/model"
	"github.com/go-git/go-git/v5"
)

var (
	urlNameReg = "^[A-Za-z0-9-_]{1,32}$"
)

//validateURLName checks if a name is legal (to be in URL)
func validateURLName(name string) bool {
	ok, _ := regexp.Match(urlNameReg, []byte(name))
	return ok
}

//GetOrgProjByName get an org project info
func GetOrgProjByName(ctx *gin.Context, owner string, name string) *model.Project {
	//TODO
	return nil
}

//GetProj get a project info
func GetProj(ctx *gin.Context) {
	uri := ctx.Param("uri")
	proj := model.GetProjByURI(uri)
	if proj == nil {
		ctx.JSON(http.StatusNotFound, model.Resp404)
		return
	}
	ctx.JSON(http.StatusOK, model.Response{
		Success: true,
		Data: gin.H{
			"proj": proj,
		},
	})
}

//ListUserProj list all projects
func ListUserProj(ctx *gin.Context) {
	user_id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	projs := model.ListUserProj(user_id)
	ctx.JSON(http.StatusOK, model.Response{
		Success: true,
		Msg:     "",
		Data: gin.H{
			"projs": projs,
		}})
}

//CreateUserProj create a new user project
func CreateUserProj(ctx *gin.Context) {
	var (
		req                model.CreateProjRequest
		srcLangs, trnLangs []model.Language
		ok                 bool
		user               = ctx.Keys["user"].(*model.User)
	)
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.Resp400)
		return
	}
	if srcLangs, ok = model.ParseLangsFromCodes(req.SrcLangs); !ok {
		ctx.JSON(http.StatusOK, model.Response{
			Success: false,
			Msg:     "src_langs illegal",
			Code:    constant.ErrProjSrcLangIllegal,
		})
		return
	}
	if trnLangs, ok = model.ParseLangsFromCodes(req.TrnLangs); !ok {
		ctx.JSON(http.StatusOK, model.Response{
			Success: false,
			Msg:     "trn_langs illegal",
			Code:    constant.ErrProjTrnLangIllegal,
		})
		return
	}
	if req.URI == "" || model.GetProjByURI(req.URI) != nil {
		ctx.JSON(http.StatusOK, model.Response{
			Success: false,
			Msg:     "uri exists",
			Code:    constant.ErrProjUriExists,
		})
		return
	}
	proj := &model.Project{
		URI:                req.URI,
		Name:               req.Name,
		OwnerID:            user.ID,
		Type:               req.Type,
		Status:             constant.ProjStatCreated,
		Desc:               req.Desc,
		GitURL:             req.GitURL,
		Path:               filepath.Join(config.ProjPath, req.URI),
		SrcLangs:           strings.Join(req.SrcLangs, constant.Delim),
		TrnLangs:           strings.Join(req.TrnLangs, constant.Delim),
		SourceLanguages:    srcLangs,
		TranslateLanguages: trnLangs,
	}
	if req.Type == constant.ProjTypeGitURL {
		//Nothing to do
	} else if req.Type == constant.ProjTypeGithub {
		proj.Token = user.GithubRepoToken
	} else {
		ctx.JSON(http.StatusOK, model.Response{
			Success: false,
			Msg:     "type illegal",
			Code:    constant.ErrProjTypeIllegal,
		})
		return
	}
	if err := proj.Create(); err != nil {
		ctx.JSON(http.StatusOK, model.Response{
			Success: false,
			Code:    constant.ErrUnknown,
			Msg:     err.Error(),
		})
	}
	ctx.JSON(http.StatusCreated, model.Response{
		Success: true,
		Msg:     "create project successfully",
	})
	go proj.Init()
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
	ctx.JSON(http.StatusOK, model.Response{
		Success: true,
		Data: gin.H{
			"branches": brchs,
		},
	})
}
