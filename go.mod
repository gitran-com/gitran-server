module github.com/wzru/gitran-server

go 1.13

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.6.3
	github.com/go-co-op/gocron v0.3.3
	github.com/go-git/go-git/v5 v5.2.0
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.3 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/sirupsen/logrus v1.7.0
	github.com/whilp/git-urls v1.0.0
	gopkg.in/yaml.v2 v2.3.0
	gorm.io/driver/mysql v1.0.3
	gorm.io/driver/postgres v1.0.5
	gorm.io/gorm v1.20.6
)
