package main

import (
	"fmt"
	"log"
	"net/http"
	"os/user"
	"path"
	"os"
	"errors"
	workers "dawidurbanski.pl/itb_uploader/workers"
)

var cPort int = 8080

var logger log.Logger


func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func FindConfigFileName() (string, error) {

	//user home directory
	usr, err := user.Current()
	if err == nil {
		fn := path.Join(usr.HomeDir, "itb_uploader.json")
		if fileExists(fn) {
			fmt.Println("configuration file found in: ", fn)
			return fn, nil
		} else {
			fmt.Println("configuration file not found in: ", fn)
		}
	}

	//curtent directory
	wd, _ := os.Getwd()
	fn := path.Join(wd, "config.json")
	if fileExists(fn) {
		fmt.Println("configuration file found in: ", fn)
		return fn, nil
	}
	fmt.Println("configuration file not found in: ", fn)

	return "", errors.New("no configuration file found")
}


func serverInit() {
	fmt.Println("starting Init")

	configFileName, err := FindConfigFileName()
	if err != nil {
		panic(err)
	}

	err = workers.LoadConfiguration(configFileName)
	if err != nil {
		panic(err)
	}
}

func startServer() {

	http.HandleFunc("/", ViewIndex)
	http.HandleFunc("/get", ViewNumber)
	http.HandleFunc("/action/date", ViewDate)
	http.HandleFunc("/action/voice", ViewVoice)
	http.HandleFunc("/action/desk/up", ViewDeskUp)
	http.HandleFunc("/action/desk/down", ViewDeskDown)
	http.HandleFunc("/upload", ViewUploadFile)

	fmt.Println("starting server")
	http.ListenAndServe(fmt.Sprintf(":%d", cPort), nil)
}

func main() {
	initLogger()
	serverInit()
	startServer()
}
