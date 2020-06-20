package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"path"
	"strconv"
	"strings"

	"github.com/UrbanskiDawid/itb_uploader/actions"
	"github.com/UrbanskiDawid/itb_uploader/logging"
	"github.com/UrbanskiDawid/itb_uploader/views"
	"github.com/spf13/cobra"
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

	err = actions.LoadConfiguration(configFileName)
	if err != nil {
		panic(err)
	}
}

func generateUserVisibleActionName(name string) string {
	var ret string
	ret = name
	ret = strings.ToLower(ret)
	ret = strings.ReplaceAll(ret, " ", "_")
	return ret
}

func startServer() {

	views.Init()

	http.HandleFunc("/", views.ViewIndex)
	http.HandleFunc("/action/", views.ViewIndex)

	for _, name := range actions.GetActionNames() {

		var actionName string
		actionName = name // note must make a copy

		var userVisibleNameName string
		userVisibleNameName = generateUserVisibleActionName(name)
		http.HandleFunc("/action/"+userVisibleNameName, views.BuildViewAction(userVisibleNameName, actionName))
	}

	fmt.Println("starting server port", port)
	logging.Log.Println("starting server port", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func main() {

	logging.InitLogger()
	configInit()

	var rootCmd = &cobra.Command{Use: "app"}
	var server = &cobra.Command{
		Use:   "server [start server]",
		Short: "server port",
		Long:  `run this app in server mode.`,
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				port, _ = strconv.ParseUint(args[0], 10, 64)
			}
			startServer()
			os.Exit(1)
		},
	}

	names := actions.GetActionNames()
	for i := 0; i < len(names); i++ {

		var name = names[i]

		file := actions.GetTargetFileNameForAction(name)

		if file != "" {
			var cmd = &cobra.Command{
				Use:   name,
				Short: fmt.Sprintf("%s [file]", name),
				Args:  cobra.ExactArgs(1),
				Run: func(cmd *cobra.Command, args []string) {

					err := actions.UploadFile(name, args[1])
					if err != nil {
						print(err)
						os.Exit(1)
					}

					stdOut, stdErr, err := actions.ExecuteAction(name)
					print(stdOut)
					if err != nil {
						print(stdErr)
						os.Exit(1)
					}
				},
			}
			rootCmd.AddCommand(cmd)
		} else {
			var cmd = &cobra.Command{
				Use:  name,
				Args: cobra.NoArgs,
				Run: func(cmd *cobra.Command, args []string) {
					stdOut, stdErr, err := actions.ExecuteAction(name)
					print(stdOut)
					if err != nil {
						print(stdErr)
						os.Exit(1)
					}
				},
			}
			rootCmd.AddCommand(cmd)
		}
	}

	//rootCmd.PersistentFlags().StringVar(&configFileName, "config", "", "Author name for copyright attribution")

	rootCmd.AddCommand(server)
	rootCmd.Execute()
}
