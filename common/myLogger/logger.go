package myLogger

import (
	"fmt"
	"ginApi/common/config"
	"log"
)

var Logid string
var mylog *log.Logger
var env bool = false // 是否为prod环境

func init() {
	mylog = log.New(LogWrite, "", log.LstdFlags)
	if config.Viper.Get("env") == "prod" {
		env = true
	}
}

func Printf(msg string, v ...any) {
	if v != nil {
		msg = fmt.Sprintf(msg, v)
	}
	if !env {
		fmt.Println("logid["+Logid+"] ", msg)
	}
	mylog.Println("logid["+Logid+"] ", msg)
}

func Println(msg string) {
	if !env {
		// 非prod环境，同时将日志输出到console
		fmt.Println("logid["+Logid+"] ", msg)
	}
	mylog.Println("logid["+Logid+"] ", msg)
}
