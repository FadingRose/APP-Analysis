package logger

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Fatal   *log.Logger
)

func init() {
	logPath := filepath.Join("../", "logs")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		os.Mkdir(logPath, os.ModePerm)
		panic(err)
	}

	fileName := time.Now().Format("2006-01-02_15-04-05") + ".log"
	filePath := filepath.Join(logPath, fileName)

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error creating log files:", err)
		panic(err)
	}

	Info = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Fatal = log.New(file, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)
}
