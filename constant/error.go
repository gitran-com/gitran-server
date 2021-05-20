package constant

const (
	//ErrorLoginOrPasswordIncorrect 用户名或密码错误
	ErrorLoginOrPasswordIncorrect = 1000 + iota
	//ErrorEmailExists 邮箱已存在
	ErrorEmailExists
)

const (
	//ErrorProjNameIllegal project名字不合法
	ErrorProjNameIllegal = 2000 + iota
	//ErrorProjSrcLangEmpty project src_langs is null
	ErrorProjSrcLangEmpty
	//ErrorProjSrcLangIllegal project src_langs is illegal
	ErrorProjSrcLangIllegal
	//ErrorProjTrnLangIllegal project trn_langs is illegal
	ErrorProjTrnLangIllegal
	//ErrorProjGitURLIllegal project git url is illegal
	ErrorProjGitURLIllegal
	//ErrorGithubUnauthorized import github project but not auth
	ErrorGithubUnauthorized
	//ErrorProjNameExists means project name exists
	ErrorProjNameExists
)

const (
	//GitErrorUpToDate git already up-to-date
	GitErrorUpToDate = "already up-to-date"
)
