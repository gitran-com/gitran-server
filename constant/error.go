package constant

const (
	//ErrLoginOrPasswordIncorrect 用户名或密码错误
	ErrLoginOrPasswordIncorrect = 1000 + iota
	//ErrEmailExists 邮箱已存在
	ErrEmailExists
)

const (
	ErrProjNameIllegal = 2000 + iota
	ErrProjUriIllegal
	ErrProjSrcLangEmpty
	ErrProjSrcLangIllegal
	ErrProjTrnLangIllegal
	ErrGitUrlIllegal
	ErrTokenIllegal
)

const (
	ErrGitChkout = 3000 + iota
	ErrGitCommit
)

const (
	//ErrGitUpToDate git already up-to-date
	ErrGitUpToDate = "already up-to-date"
)
