package logger

import (
	"log"
	"os"
)

func InitLogger(filepath string) *log.Logger {
	if filepath == "" {
		return nil
	}

	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	return log.New(f, "", log.LstdFlags)
}
