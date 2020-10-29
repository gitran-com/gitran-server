package constant

import "github.com/gin-gonic/gin"

const (
	DebugMode   = gin.DebugMode
	ReleaseMode = gin.ReleaseMode
	TestMode    = gin.TestMode
)

const (
	//OwnerUsr 属于个人的项目
	OwnerUsr = iota
	//OwnerOrg 属于组织的项目
	OwnerOrg
)

const (
	//ProjCmn 普通项目
	ProjCmn = iota
	//ProjGit git仓库引入的项目
	ProjGit
)

const (
	//StatusCreated 已创建未初始化
	StatusCreated = iota
	//StatusInit 已初始化完成
	StatusInit
)

const (
	//Delimiter 是数据库记录中的分隔符
	Delimiter = "|"
)
