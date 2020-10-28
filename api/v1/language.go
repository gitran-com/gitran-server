package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/wzru/gitran-server/controller"
)

func langInit(g *gin.RouterGroup) {
	gg := g.Group("/languages")
	gg.GET("", controller.GetLangs)
	gg.GET("/:id", controller.GetLang)
}
