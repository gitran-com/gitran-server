package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/wzru/gitran-server/constant"
)

type oauth struct {
	Enable       bool   `json:"enable"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type db struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	User        string `json:"user"`
	Password    string `json:"password"`
	Host        string `json:"host"`
	Port        uint16 `json:"port"`
	TablePrefix string `json:"table_prefix"`
}

type jwt struct {
	Secret       string   `json:"secret"`
	SkipperPaths []string `json:"skipper_paths"`
	ValidTime    uint     `json:"valid_time"`
	RefreshTime  uint     `json:"refresh_time"`
}

type app struct {
	Name          string `json:"name"`
	URL           string `json:"url"`
	Addr          string `json:"addr"`
	APIPrefix     string `json:"api_prefix"`
	SessionSecret string `json:"session_secret"`
}

type email struct {
	Enable   bool   `json:"enable"`
	From     string `json:"from"`
	Type     string `json:"type"`
	Host     string `json:"host"`
	TLS      bool   `json:"tls"`
	User     string `json:"user"`
	Password string `json:"password"`
}

//Config 应用运行配置
type Config struct {
	Github    oauth    `json:"github"`
	DB        db       `json:"db"`
	JWT       jwt      `json:"jwt"`
	APP       app      `json:"app"`
	Email     email    `json:"email"`
	FileTypes []string `json:"file_types"`
}

type logFormatter struct{}

var (
	//Mode 运行模式
	Mode = flag.String("mode", constant.DebugMode, "运行模式")
	//IsDebug 是否是调试模式
	IsDebug    = false
	configFile = flag.String("config", "config.json", "配置文件路径")
	pwd, _     = os.Getwd()
	logPath    = pwd + "/log/"
	//DataPath 数据目录
	DataPath = "data/"
	//ProjPath 项目目录
	ProjPath    = DataPath + "project/"
	mainLogFile = logPath + "gitran.log"
	//GinLogFile Gin日志目录
	GinLogFile = logPath + "api"
	//TimeFormat 日志时间格式
	TimeFormat = "2006/01/02 15:04:05"
)

var (
	c *Config
	//DB config
	DB *db
	//Github config
	Github *oauth
	//JWT config
	JWT *jwt
	//APP config
	APP *app
	//Email config
	Email *email
)

func (s *logFormatter) Format(entry *log.Entry) ([]byte, error) {
	msg := fmt.Sprintf("%s [%s] %s\n", time.Now().Local().Format(TimeFormat), strings.ToUpper(entry.Level.String()), entry.Message)
	if IsDebug {
		fmt.Printf("%s", msg)
	}
	return []byte(msg), nil
}

//Init init config
func Init() error {
	gin.SetMode(*Mode)
	IsDebug = (*Mode == constant.DebugMode)
	//打开日志文件
	logFile, err := os.OpenFile(mainLogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("Cannot open main log file '%s'! Try create the directory.\n", mainLogFile)
		return err
	}
	log.SetFormatter(new(logFormatter))
	log.SetOutput(logFile)
	//加载配置文件
	log.Infof("open config file '%s'...", *configFile)
	configData, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("cannot open config file '%s'!\n", *configFile)
		return err
	}
	// fmt.Printf("%v\n", json.Valid(configData))
	c = &Config{}
	if err := json.Unmarshal(configData, c); err != nil {
		log.Fatalf("config JSON unmarshal failed: %v", err)
		return err
	}
	Github = &c.Github
	DB = &c.DB
	JWT = &c.JWT
	APP = &c.APP
	Email = &c.Email
	return nil
}
