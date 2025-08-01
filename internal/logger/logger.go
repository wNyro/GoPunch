package logger

import (
	"log"
	"os"
)

func InitLogger(filepath string) *log.Logger {
	if filepath == "" {
		return log.New(os.Stdout, "", log.LstdFlags)
	}

	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Failed to open log file: %v, falling back to stdout", err)
		return log.New(os.Stdout, "", log.LstdFlags)
	}
	return log.New(f, "", log.LstdFlags|log.Lshortfile)
}
