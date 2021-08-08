package constant

const (
	ErrUnknown = -1
)

// auth error
const (
	//ErrEmailOrPassIncorrect 邮箱或密码错误
	ErrEmailOrPassIncorrect = 1000 + iota
	//ErrEmailExists 邮箱已存在
	ErrEmailExists
	ErrGithubRepoUnauthorized
)

const (
	ErrProjUriExists = 2000 + iota
	ErrProjSrcLangEmpty
	ErrProjSrcLangIllegal
	ErrProjTrnLangIllegal
	ErrProjTypeIllegal
)

const (
	ErrGitChkout = 3000 + iota
	ErrGitCommit
)

const (
	ErrTokenExpired = 4000 + iota
)

const (
	//ErrGitUpToDate git already up-to-date
	ErrGitUpToDate = "already up-to-date"
)
