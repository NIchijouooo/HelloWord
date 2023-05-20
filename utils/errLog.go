package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

var ErrLog *log.Logger

func ErrorLogInit() {

	file := "./log/" + time.Now().Format("2006-01-02") + "_errlog" + ".txt"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		fmt.Printf("openLogFile err,%v", err)
		return
	}
	// 将文件设置为loger作为输出
	ErrLog = log.New(logFile, "", log.LstdFlags|log.Lshortfile|log.LUTC)
}
