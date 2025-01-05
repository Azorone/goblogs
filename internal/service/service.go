package service

import (
	"log"
	"mastcat/internal/config"
	"os"
)

func InitLogFile() {
	logFilePath := config.AppConfig.LogFile
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		file, err := os.Create(logFilePath)
		if err != nil {
			log.Fatalf("Failed to create log file: %v", err)
		}
		file.Close()
	}
}
