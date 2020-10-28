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

type database struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	User        string `json:"user"`
	Password    string `json:"password"`
	Host        string `json:"host"`
	TablePrefix string `json:"table_prefix"`
}

type jwt struct {
	Secret       string   `json:"secret"`
	SkipperPaths []string `json:"skipper_paths"`
	ValidTime    int64    `json:"valid_time"`
	RefreshTime  int64    `json:"refresh_time"`
}

type app struct {
	Name      string `json:"name"`
	Addr      string `json:"addr"`
	APIPrefix string `json:"api_prefix"`
}

type Lang struct {
	ID    uint   `json:"id"`
	Code  string `json:"code"`
	Code3 string `json:"code3"`
	ISO   string `json:"iso"`
	Name  string `json:"name"`
}

type Config struct {
	Github    oauth    `json:"github"`
	DB        database `json:"database"`
	JWT       jwt      `json:"jwt"`
	APP       app      `json:"app"`
	Langs     []Lang   `json:"langs"`
	FileTypes []string `json:"file_types"`
}

type logFormatter struct{}

var (
	Mode        = flag.String("mode", constant.DebugMode, "运行模式")
	IsDebug     = false
	configFile  = flag.String("config", "config.json", "配置文件路径")
	pwd, _      = os.Getwd()
	logPath     = pwd + "/log/"
	DataPath    = pwd + "/data/"
	mainLogFile = logPath + "gitran.log"
	GinLogFile  = logPath + "api"
	TimeFormat  = "2006/01/02 15:04:05"
)

var (
	c *Config
	//DB config
	DB *database
	//Github config
	Github *oauth
	//JWT config
	JWT *jwt
	//APP config
	APP *app
	//Langs config
	Langs []Lang
)

func (s *logFormatter) Format(entry *log.Entry) ([]byte, error) {
	msg := fmt.Sprintf("%s [%s] %s\n", time.Now().Local().Format(TimeFormat), strings.ToUpper(entry.Level.String()), entry.Message)
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
	log.Infof("Open config file '%s'...", *configFile)
	configData, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("Cannot open config file '%s'!\n", *configFile)
		return err
	}
	// fmt.Printf("%v\n", json.Valid(configData))
	c = &Config{}
	if err := json.Unmarshal(configData, c); err != nil {
		log.Fatalf("Config JSON unmarshal failed: %v", err)
		return err
	}
	// fmt.Println("hello")
	Github = &c.Github
	DB = &c.DB
	JWT = &c.JWT
	APP = &c.APP
	Langs = c.Langs
	//fmt.Printf("%v", *C)
	return nil
}
