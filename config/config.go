package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	ConfigFile = flag.String("config", "config.json", "配置文件路径")
	LogFile    = flag.String("log", "log/gitran.log", "日志文件路径")
)

type Oauth struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	CallbackURL  string `json:"callback_url"`
}

type Database struct {
	Type   string `json:"type"`
	Source string `json:"source"`
}

type LogFormatter struct{}

func (s *LogFormatter) Format(entry *log.Entry) ([]byte, error) {
	msg := fmt.Sprintf("%s [%s] %s\n", time.Now().Local().Format("2006/01/02 15:04:05"), strings.ToUpper(entry.Level.String()), entry.Message)
	return []byte(msg), nil
}

func Init() error {
	logFile, err := os.OpenFile(*LogFile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Cannot open log file '%s'! Try new the directory")
		return err
	}
	log.SetFormatter(new(LogFormatter))
	log.SetOutput(logFile)
	log.Infof("Open config file '%s'...", *ConfigFile)
	configJSON, err := os.Open(*ConfigFile)
	defer configJSON.Close()
	if err != nil {
		log.Fatalf("Cannot open config file '%s'!\n", *ConfigFile)
		return err
	}
	return nil
}
