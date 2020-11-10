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
	OwnerUsr uint8 = iota
	//OwnerOrg 属于组织的项目
	OwnerOrg
)

const (
	//TypeCommon 普通项目
	TypeCommon uint8 = iota
	//TypeGithub GitHub引入的项目
	TypeGithub
	//TypeGitURL Git URL引入的公开项目
	TypeGitURL
)

const (
	//ProjStatCreated 项目已创建未初始化
	ProjStatCreated uint8 = iota
	//ProjStatInit 项目已初始化完成
	ProjStatInit
)

const (
	//RuleStatCreated 规则已创建未保存到配置文件
	RuleStatCreated uint8 = iota
	//RuleStatSaved 规则已保存到配置文件
	RuleStatSaved
)

const (
	//SyncStatSucc 上一次分支同步成功
	SyncStatSucc uint8 = iota
	//SyncStatDoing 分支同步中
	SyncStatDoing
	//SyncStatFail 上一次分支同步失败
	SyncStatFail
)

const (
	//AuthGithubLogin github登录
	AuthGithubLogin = "login"
	//AuthGithubImport github引入仓库
	AuthGithubImport = "import"
)

const (
	//Delim 是数据库记录中的分隔符
	Delim = "|"
	//GithubTokenScopeDelim 是GitHub token作用域的分隔符
	GithubTokenScopeDelim = ","
)
