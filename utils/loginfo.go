package utils

import (
	"log"
	"os"
)

const (
	//定义日志级别
	UNKNOWN = iota
	DEBUG
	TRACE
	INFO
	WARNING
	ERROR
	FATAL
)

// logName表示要输出的日志文件名
func Log(logFileName string, level int) *log.Logger {
	var (
		logFile *os.File
		err     error
	)
	var levelString string
	switch level {
	case UNKNOWN:
		levelString = "[UNKNOWN] "
	case DEBUG:
		levelString = "[DEBUG] "
	case TRACE:
		levelString = "[TRACE] "
	case INFO:
		levelString = "[INFO] "
	case WARNING:
		levelString = "[WARNING] "
	case ERROR:
		levelString = "[ERROR] "
	case FATAL:
		levelString = "[FATAL] "
	}
	if len(logFileName) == 0 {
		logFile = os.Stderr
	} else {
		str, _ := os.Getwd()
		logFile, err = os.OpenFile(str+"/log/"+logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	}
	if err != nil {
		log.Fatal("日志文件无法正常打开", err)
	}
	return log.New(logFile, levelString, log.Ldate|log.Ltime|log.Lshortfile)
}
