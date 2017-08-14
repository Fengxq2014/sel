package tool

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Info Info日志
func Info(v ...interface{}) {
	writeLog("INFO ", v)
}

// Debug debug日志
func Debug(v ...interface{}) {
	writeLog("DEBUG ", v)
}

// Error error日志
func Error(v ...interface{}) {
	writeLog("ERROR ", v)
}

func getFileName() string {
	pwd, _ := os.Getwd()
	fileName := "./" + time.Now().Format("2006-01-02") + ".log"
	return filepath.Join(pwd, "log", fileName)
}

func writeLog(level string, v ...interface{}) {
	myfile, _ := os.OpenFile(getFileName(), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	logger := log.New(io.MultiWriter(myfile, os.Stdout), level, log.Ldate|log.Ltime|log.Llongfile)
	logger.Output(3, fmt.Sprintln(v))
}
