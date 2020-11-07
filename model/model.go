package model

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"github.com/wzru/gitran-server/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db       *gorm.DB
	langs    []Language
	langFile = flag.String("language", "language.json", "语言列表文件路径")
)

func initDB() error {
	var err error
	if config.DB.Type == "mysql" {
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8mb4&parseTime=True&loc=Local", config.DB.User, config.DB.Password, config.DB.Host)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Database connect ERROR : %v", err)
			return err
		}
		db.Exec("CREATE DATABASE IF NOT EXISTS " + config.DB.Name)
		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DB.User, config.DB.Password, config.DB.Host, config.DB.Name)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	} else if config.DB.Type == "postgresql" {
		//TODO
	}
	if err != nil {
		log.Fatalf("DB connect ERROR : %v", err.Error())
		return err
	}
	err = db.AutoMigrate(&User{}, &Project{}, &ProjCfg{}, &BrchRule{}, &Translation{})
	if err != nil {
		log.Fatalf("DB migrate ERROR : %v", err.Error())
		return err
	}
	return nil
}

func initLangs() error {
	langData, err := ioutil.ReadFile(*langFile)
	if err != nil {
		log.Fatalf("cannot open language file '%s'!\n", *langFile)
		return err
	}
	if err := json.Unmarshal(langData, &langs); err != nil {
		log.Fatalf("language JSON unmarshal failed: %v", err)
		return err
	}
	// fmt.Printf("langs:%+v", langs)
	return nil
}

//Init initialize the model
func Init() error {
	if err := initDB(); err != nil {
		return err
	}
	return initLangs()
}
