package logger

import (
	"fmt"
	"time"
)

// Test Log

func Debug(format string,args ...interface{})  {
	date := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(fmt.Sprintf("[%s] [Debug] %s",date,fmt.Sprintf(format,args...)))
}
func Error(format string,args ...interface{})  {
	date := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(fmt.Sprintf("[%s] [Error] %s",date,fmt.Sprintf(format,args...)))
}
func Info(format string,args ...interface{})  {
	date := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(fmt.Sprintf("[%s] [Info] %s",date,fmt.Sprintf(format,args...)))
}