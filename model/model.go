package model

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/gitran-com/gitran-server/config"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db       *gorm.DB
	langs    []Language
	langFile = flag.String("language", "language.json", "语言列表文件路径")
	langMap  = make(map[string]Language)
)

func initDB() error {
	var err error
	var dsn string
	switch config.DB.Type {
	case "mysql", "mariadb":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%v)/?charset=utf8mb4&parseTime=True&loc=Local", config.DB.User, config.DB.Password, config.DB.Host, config.DB.Port)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("db connect ERROR : %v", err)
			return err
		}
		db.Exec("CREATE DATABASE IF NOT EXISTS " + config.DB.Name)
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DB.User, config.DB.Password, config.DB.Host, config.DB.Port, config.DB.Name)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "postgres", "postgresql":
		dsn = fmt.Sprintf("host=%s port=%v dbname=postgres user=%v password='%s' sslmode=disable", config.DB.Host, config.DB.Port, config.DB.User, config.DB.Password)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("db connect ERROR : %v", err)
			return err
		}
		db = db.Exec("CREATE DATABASE " + config.DB.Name)
		if db.Error != nil {
			log.Warnf("db create %v ERROR : %v", config.DB.Name, db.Error.Error())
			// return db.Error
		}
		dsn = fmt.Sprintf("host=%s port=%v dbname=%v user=%v password='%s' sslmode=disable", config.DB.Host, config.DB.Port, config.DB.Name, config.DB.User, config.DB.Password)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	default:
		log.Fatalf("unknown db type : %v", config.DB.Type)
	}
	if err != nil {
		log.Fatalf("db connect ERROR : %v", err.Error())
		return err
	}
	err = db.AutoMigrate(&User{}, &Project{}, &ProjCfg{}, &ProjRole{}, &ProjFile{}, &Sentence{}, &Translation{})
	if err != nil {
		log.Fatalf("db migrate ERROR : %v", err.Error())
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
	for _, lang := range langs {
		langMap[lang.Code] = lang
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
