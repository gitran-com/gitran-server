package constant

import "github.com/gin-gonic/gin"

const (
	DebugMode   = gin.DebugMode
	ReleaseMode = gin.ReleaseMode
	TestMode    = gin.TestMode
)

const (
	//OwnerUsr 属于个人的项目
	OwnerUsr = 0
	//OwnerOrg 属于组织的项目
	OwnerOrg = 1
)

const (
	//ProjCmn 普通项目
	ProjCmn = 0
	//ProjGit git仓库引入的项目
	ProjGit = 1
)

const (
	//Delimiter 是数据库记录中的分隔符
	Delimiter = "|"
)
