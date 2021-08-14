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
	OwnerUsr int = iota
	//OwnerOrg 属于组织的项目
	OwnerOrg
)

const (
	//ProjTypePlain 普通项目
	ProjTypePlain int = iota
	//ProjTypeGitURL Git URL引入的公开项目
	ProjTypeGitURL
	//ProjTypeGithub Github token引入的项目
	ProjTypeGithub
)

const (
	MaxProjInitRetry        = 5
	MinProjPullGap   uint16 = 5
	MinProjPushGap   uint16 = 5
)

const (
	//ProjStatCreated 项目已创建未初始化
	ProjStatCreated int = iota
	//ProjStatReady 项目已就绪
	ProjStatReady
	//ProjStatProcessingString 项目正在处理string
	ProjStatProcessingString
	//ProjStatFailed 项目初始化失败
	ProjStatFailed int = -1
)

const (
	//RuleStatCreated 规则已创建未保存到配置文件
	RuleStatCreated int = iota
	//RuleStatSaved 规则已保存到配置文件
	RuleStatSaved
)

const (
	//SyncStatSucc 上一次分支同步成功
	SyncStatSucc int = iota
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

const (
	SubjLogin            = "login"
	SubjRegister         = "register"
	SubjRefresh          = "refresh"
	SubjGithubLogin      = "github-login"
	SubjGithubFirstLogin = "github-first-login"
)
