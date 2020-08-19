package model

import (
	"github.com/WangZhengru/gitran-be/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

var DB *gorm.DB

func Init() error {
	DB, err := gorm.Open(config.C.Database.Type, config.C.Database.Source)
	if err != nil {
		log.Fatalf("Database open ERROR : %v", err.Error())
		return err
	}
	DB.Exec("CREATE DATABASE IF NOT EXISTS " + config.C.Database.Name)
	DB.Exec("USE " + config.C.Database.Name)
	DB.AutoMigrate(&User{}, &Project{}, &Translation{})
	return nil
}
