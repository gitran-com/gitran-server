package constant

import "github.com/gin-gonic/gin"

const (
	DebugMode   = gin.DebugMode
	ReleaseMode = gin.ReleaseMode
	TestMode    = gin.TestMode
)

const (
	ProjGen = "gen" //普通项目
	ProjGit = "git" //从git仓库引入的项目
)

const (
	Delimiter = "|"
)
