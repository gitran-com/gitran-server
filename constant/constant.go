package constant

import "github.com/gin-gonic/gin"

const (
	//DebugMode 调试模式
	DebugMode = gin.DebugMode
	//ReleaseMode 发布模式
	ReleaseMode = gin.ReleaseMode
	//TestMode 测试模式
	TestMode = gin.TestMode
)

const (
	//OwnerUsr 属于个人的项目
	OwnerUsr = iota
	//OwnerOrg 属于组织的项目
	OwnerOrg
)

const (
	//ProjCommon 普通项目
	ProjCommon = iota
	//ProjGithub git仓库引入的项目
	ProjGithub
)

const (
	//ProjStatCreated 项目已创建未初始化
	ProjStatCreated = iota
	//ProjStatInit 项目已初始化完成
	ProjStatInit
)

const (
	//Delim 是数据库记录中的分隔符
	Delim = "|"
)
