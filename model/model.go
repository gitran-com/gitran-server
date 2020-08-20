package model

import (
	"github.com/WangZhengru/gitran-be/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

var DB *gorm.DB

func Init() error {
	var err error
	DB, err = gorm.Open(config.DB.Type, config.DB.Source)
	if err != nil {
		log.Fatalf("Database open ERROR : %v", err.Error())
		return err
	}
	DB.Exec("CREATE DATABASE IF NOT EXISTS " + config.DB.Name)
	DB.Exec("USE " + config.DB.Name)
	DB.AutoMigrate(&User{}, &Project{}, &Translation{})
	return nil
}
