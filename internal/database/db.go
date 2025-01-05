package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"mastcat/internal/config"
	"mastcat/internal/model"
	"os"
)

var DB *gorm.DB
var err error

func InitDB() {
	//dsn := "root:WoDeMySQLMiMaShi123456YiDingYaoJiZhu@tcp(127.0.0.1:64521)/catmastcc?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBName,
	)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}
	err = DB.AutoMigrate(&model.User{}, &model.Blog{}, &model.Category{}, &model.Image{}, &model.AccessLog{})
	if err != nil {
		os.Exit(1)
	}
}
