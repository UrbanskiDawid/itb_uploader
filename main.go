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

	"github.com/UrbanskiDawid/itb_uploader/actions"
	"github.com/UrbanskiDawid/itb_uploader/logging"
	"github.com/UrbanskiDawid/itb_uploader/views"
	"github.com/urfave/cli"
)

var logger log.Logger

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
		return fn, nil
	}
	return "", errors.New("no configuration file found")
}

func configInit() {
	fmt.Println("config Init")

	configFileName, err := findConfigFileName()
	if err != nil {
		panic(err)
	}

	err = actions.LoadConfiguration(configFileName)
	if err != nil {
		panic(err)
	}
	fmt.Println("logging to: " + configFileName)
}

func generateUserVisibleActionName(name string) string {
	var ret string
	ret = name
	ret = strings.ToLower(ret)
	ret = strings.ReplaceAll(ret, " ", "_")
	return ret
}

func startServer(port int) {

	views.Init()

	http.HandleFunc("/", views.ViewIndex)
	http.HandleFunc("/get", views.ViewNumber)
	http.HandleFunc("/upload", views.ViewUploadFile)

	http.HandleFunc("/action/", views.ViewAllActions)

	for _, name := range actions.GetActionNames() {

		var actionName string
		actionName = name // note must make a copy

		var userVisibleNameName string
		userVisibleNameName = generateUserVisibleActionName(name)
		http.HandleFunc("/action/"+userVisibleNameName, views.BuildViewAction(userVisibleNameName, actionName))
		println("/action/" + userVisibleNameName)
	}

	fmt.Println("starting server port", port)
	logging.Log.Println("starting server port", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

//doc: https://github.com/urfave/cli/blob/master/docs/v1/manual.md
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

	actionNames := actions.GetActionNames()
	for _, name := range actionNames {

		var actionName string
		actionName = name // note must make a copy

		var userVisibleNameName string
		userVisibleNameName = generateUserVisibleActionName(name)

		cmd := cli.Command{
			Name: userVisibleNameName,
			Action: func(c *cli.Context) error {
				println("staring action '", actionName, "'")
				stdOut, stdErr, err := actions.ExecuteAction(actionName)
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
