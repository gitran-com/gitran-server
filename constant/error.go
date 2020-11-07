package constant

const (
	//ErrorLoginOrPasswordIncorrect 用户名或密码错误
	ErrorLoginOrPasswordIncorrect = 1000 + iota
	//ErrorLoginExists 用户名已存在
	ErrorLoginExists
	//ErrorEmailExists 邮箱已存在
	ErrorEmailExists
)
