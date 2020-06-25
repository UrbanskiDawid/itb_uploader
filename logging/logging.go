package logging

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Logger *log.Logger
)

func LogConsole(logMsg string) {
	//errLog..Write(logMsg)
	Logger.Print(logMsg)
	fmt.Println(logMsg)
}

func InitLogger(logFileName string) {

	e, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Printf("error opening file: %v", err)
		os.Exit(1)
	}
	Logger = log.New(e, "", log.Ldate|log.Ltime)

	Logger.SetOutput(&lumberjack.Logger{
		Filename:   logFileName,
		MaxSize:    1,  // megabytes after which new file is created
		MaxBackups: 3,  // number of backups
		MaxAge:     28, //days
	})

}
