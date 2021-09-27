package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gitran-com/gitran-server/model"
)

func AddLikes(ctx *gin.Context) {
	var (
		tran_id, _ = strconv.ParseInt(ctx.Param("tran_id"), 10, 64)
		user       = ctx.Keys["user"].(*model.User)
	)
	ctx.JSON(http.StatusCreated, model.Response{
		Success: true,
		Data: gin.H{
			"likes": model.TxnVote(user.ID, tran_id, model.VoteLike),
		},
	})
}

func AddUnlikes(ctx *gin.Context) {
	var (
		tran_id, _ = strconv.ParseInt(ctx.Param("tran_id"), 10, 64)
		user       = ctx.Keys["user"].(*model.User)
	)
	ctx.JSON(http.StatusCreated, model.Response{
		Success: true,
		Data: gin.H{
			"likes": model.TxnVote(user.ID, tran_id, model.VoteUnlike),
		},
	})
}

func DelVote(ctx *gin.Context) {
	var (
		tran_id, _ = strconv.ParseInt(ctx.Param("tran_id"), 10, 64)
		user       = ctx.Keys["user"].(*model.User)
	)
	ctx.JSON(http.StatusOK, model.Response{
		Success: true,
		Data: gin.H{
			"likes": model.TxnVote(user.ID, tran_id, model.VoteNone),
		},
	})
}
