package logging

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	Log *log.Logger
)

func LogConsole(logMsg string) {
	Log.Println(logMsg)
	fmt.Println(logMsg)
}

func InitLogger(name string) {
	// set location of log file
	var logpath = name + ".log"

	flag.Parse()
	var file, err1 = os.Create(logpath)

	if err1 != nil {
		panic(err1)
	}
	Log = log.New(file, "", log.LstdFlags|log.Lshortfile)
	Log.Println("LogFile : " + logpath)
}
