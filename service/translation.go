package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/model"
)

func ListSentTrans(ctx *gin.Context) {
	sent_id, _ := strconv.ParseInt(ctx.Param("sent_id"), 10, 64)
	ctx.JSON(http.StatusOK, model.Response{
		Success: true,
		Data: gin.H{
			"trans": model.ListSentTrans(sent_id),
		},
	})
}

func PostTran(ctx *gin.Context) {
	var (
		req  model.PostTranRequest
		user = ctx.Keys["user"].(*model.User)
		proj = ctx.Keys["proj"].(*model.Project)
	)
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.Resp400)
		return
	}
	sent_id, _ := strconv.ParseInt(ctx.Param("sent_id"), 10, 64)
	sent := model.GetSentByID(sent_id)
	if sent == nil {
		ctx.JSON(http.StatusNotFound, model.Resp404)
		return
	}
	if sent.Locked {
		ctx.JSON(http.StatusForbidden, model.Resp403)
		return
	}
	tran := model.GetTran(sent_id, user, req.LangCode)
	if tran == nil {
		tran = &model.Translation{
			ID:       sent_id,
			UserID:   user.ID,
			UserName: user.Name,
			ProjID:   proj.ID,
			FileID:   req.FileID,
			SentID:   sent_id,
			Content:  req.Content,
			LangCode: req.LangCode,
		}
	} else {
		tran.Content = req.Content
	}
	tran.Write()
	ctx.JSON(http.StatusCreated, model.Response{
		Success: true,
		Data: gin.H{
			"tran": tran,
		},
	})
}
