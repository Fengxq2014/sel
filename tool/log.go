package tool

import (
	"io"
	"time"
	"os"
	"path/filepath"
	"log"
)

func GetLogger() *log.Logger{
	pwd, _ := os.Getwd()
	fileName := "./" + time.Now().Format("2006-01-02") + ".log"
	s := filepath.Join(pwd, "log", fileName)
	myfile, _ := os.OpenFile(s, os.O_APPEND|os.O_CREATE|os.O_RDWR, 066)
	logger := log.New(io.MultiWriter(myfile, os.Stdout),"",log.Ldate|log.Ltime|log.Llongfile)
	return logger
}