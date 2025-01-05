package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"mastcat/internal/config"
	"mastcat/internal/database"
	"mastcat/internal/router"
	"mastcat/internal/service"
	"os"
)

func loadConfig(env string) {
	file, err := os.Open("configs/config." + env + ".json")
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config.AppConfig); err != nil {
		log.Fatalf("Failed to decode config file: %v", err)
	}
}
func main() {

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev" // default to development environment
	}

	loadConfig(env)

	database.InitDB()
	service.InitLogFile()
	f, _ := os.Open(config.AppConfig.LogFile)
	gin.DefaultWriter = io.MultiWriter(f)
	r := router.SetupRouter()
	if err := r.Run(config.AppConfig.Address + ":" + config.AppConfig.Port); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
