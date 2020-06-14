package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/UrbanskiDawid/itb_uploader/logging"
	"github.com/UrbanskiDawid/itb_uploader/views"
	"github.com/UrbanskiDawid/itb_uploader/workers"
	"github.com/urfave/cli" // imports as package "cli"
)

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

func configInit() {
	fmt.Println("config Init")

	configFileName, err := FindConfigFileName()
	if err != nil {
		panic(err)
	}

	err = workers.LoadConfiguration(configFileName)
	if err != nil {
		panic(err)
	}
	fmt.Println("logging to " + configFileName)

}

func startServer(port int) {

	http.HandleFunc("/", views.ViewIndex)
	http.HandleFunc("/get", views.ViewNumber)
	http.HandleFunc("/upload", views.ViewUploadFile)

	http.HandleFunc("/action/", views.ViewAllActions)
	http.HandleFunc("/action/date", views.ViewDate)
	http.HandleFunc("/action/voice", views.ViewVoice)
	http.HandleFunc("/action/desk/up", views.ViewDeskUp)
	http.HandleFunc("/action/desk/down", views.ViewDeskDown)

	fmt.Println("starting server port", port)
	logging.Log.Println("starting server port", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func argsParse() {
	app := cli.NewApp()
	app.Usage = "make an explosive entrance"

	var configServerPort int

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "port",
			Value:       8080,
			Usage:       "language for the greeting",
			Destination: &configServerPort,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "server",
			Usage: "complete a task on the list",
			Action: func(c *cli.Context) error {
				startServer(configServerPort)
				return nil
			},
		},
	}

	actionNames := workers.GetActionNames()
	for _, name := range actionNames {

		var actionName string

		actionName = name
		actionName = strings.ToLower(actionName)
		actionName = strings.ReplaceAll(actionName, " ", "_")

		cmd := cli.Command{
			Name: actionName,
			Action: func(c *cli.Context) error {
				println("staring action '", actionName, "'")
				stdOut, stdErr, err := workers.ExecuteAction(actionName)
				print(stdOut)
				print(stdErr)
				return err
			},
		}

		app.Commands = append(app.Commands, cmd)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	logging.InitLogger()
	configInit()
	argsParse()
}
