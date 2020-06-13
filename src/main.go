package main

import (
	"fmt"
	"log"
	"net/http"
	"os/user"
	"path"
	"os"
	"errors"
	workers "github.com/UrbanskiDawid/itb_uploader/workers"
	views "github.com/UrbanskiDawid/itb_uploader/views"
	logging "github.com/UrbanskiDawid/itb_uploader/logging"
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

	http.HandleFunc("/",                views.ViewIndex)
	http.HandleFunc("/get",             views.ViewNumber)
	http.HandleFunc("/action/date",     views.ViewDate)
	http.HandleFunc("/action/voice",    views.ViewVoice)
	http.HandleFunc("/action/desk/up",  views.ViewDeskUp)
	http.HandleFunc("/action/desk/down",views.ViewDeskDown)
	http.HandleFunc("/upload",          views.ViewUploadFile)

	fmt.Println("starting server")
	logging.Log.Println("starting server")
	http.ListenAndServe(fmt.Sprintf(":%d", cPort), nil)
}

func main() {
	logging.InitLogger()
	serverInit()
	startServer()
}
