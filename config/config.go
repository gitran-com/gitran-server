package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type oauth struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	CallbackURL  string `json:"callback_url"`
}

type database struct {
	Type   string `json:"type"`
	Name   string `json:"name"`
	Source string `json:"source"`
}

type jwt struct {
	Secret       string   `json:"secret"`
	SkipperPaths []string `json:"skipper_paths"`
	ValidTime    int64    `json:"valid_time"`
}

type app struct {
	Addr      string `json:"addr"`
	APIPrefix string `json:"api_prefix"`
}

type Config struct {
	Github   oauth    `json:"github"`
	Database database `json:"database"`
	JWT      jwt      `json:"jwt"`
	APP      app      `json:"app"`
}

type logFormatter struct{}

var (
	configFile  = flag.String("config", "config.json", "配置文件路径")
	logPath     = "log/"
	mainLogFile = logPath + "gitran.log"
	GinLogFile  = logPath + "gin.log"
	C           *Config
	TimeFormat  = "2006/01/02 15:04:05"
	DataPath    = "data/"
)

func (s *logFormatter) Format(entry *log.Entry) ([]byte, error) {
	msg := fmt.Sprintf("%s [%s] %s\n", time.Now().Local().Format(TimeFormat), strings.ToUpper(entry.Level.String()), entry.Message)
	return []byte(msg), nil
}

func Init() error {
	//打开日志文件
	logFile, err := os.OpenFile(mainLogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Cannot open main log file '%s'! Try new the directory")
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
	//fmt.Printf("%v\n", json.Valid(configData))
	C = &Config{}
	if err := json.Unmarshal(configData, C); err != nil {
		log.Fatalf("Config JSON unmarshal failed: %v", err)
		return err
	}
	fmt.Printf("%v", *C)
	return nil
}
