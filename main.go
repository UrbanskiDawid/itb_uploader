package main

import (
	"errors"
	"log"
	"os"
	"os/user"
	"path"

	"github.com/UrbanskiDawid/itb_uploader/actions"
	"github.com/UrbanskiDawid/itb_uploader/logging"
)

var logger log.Logger
var port uint64 = 8080

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func findConfigFileName() (string, error) {

	//user home directory
	usr, err := user.Current()
	if err == nil {
		fn := path.Join(usr.HomeDir, "itb_uploader.json")
		if fileExists(fn) {
			logging.Log.Println("configuration file found in: ", fn)
			return fn, nil
		} else {
			logging.Log.Println("configuration file not found in: ", fn)
		}
	}

	//curtent directory
	wd, _ := os.Getwd()
	fn := path.Join(wd, "config.json")
	if fileExists(fn) {
		return fn, nil
	}
	return "", errors.New("no configuration file found")
}

func configInit() {
	logging.Log.Println("config Init")

	configFileName, err := findConfigFileName()
	if err != nil {
		panic(err)
	}

	err = actions.Init(configFileName)
	if err != nil {
		panic(err)
	}

	logging.Log.Println("config:", configFileName)
}

func main() {
	logging.InitLogger("itb_uploader")
	configInit()
	runCli()
}
