package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/UrbanskiDawid/itb_uploader/actions"
	"github.com/UrbanskiDawid/itb_uploader/views"
	"github.com/spf13/cobra"
)

func buildCommand(executor actions.Executor) *cobra.Command {

	action := executor.GetAction()
	name := strings.ReplaceAll(action.Name, " ", "_")

	var cmd = &cobra.Command{
		Use: name,
		Run: func(cmd *cobra.Command, args []string) {

			if action.HasUploadFile() {
				err, _ := executor.UploadFile(args[1])
				if err != nil {
					print(err)
					os.Exit(1)
				}
			}

			stdOut, stdErr, err := executor.Execute()
			print(stdOut)
			if err != nil {
				print(stdErr)
				os.Exit(1)
			}

			if action.HasDownloadFile() {

				targetFileName := "download"
				//logging.Log.Printf("download %s to %s", actions.GetDownloadFileNameForAction(actionName), targetFileName)

				err, _ := executor.DownloadFile("test")
				if err != nil {
					print("ERROR", err)
					print(stdErr)
					os.Exit(1)
				}
				println("new file:", targetFileName)
			}
		},
	}

	short := action.Name
	if action.HasUploadFile() {
		cmd.Args = cobra.ExactArgs(1)
		short = fmt.Sprintf("%s [file]", short)
	}
	if action.HasDownloadFile() {
		short = fmt.Sprintf("%s this cmd will download file", short)
	}
	if action.HasCommand() {
		short = fmt.Sprintf("%s this cmd will run command '%s'", short, action.Cmd)
	}
	cmd.Short = short
	return cmd
}

func runCli() {

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
			views.StartServer(port, actions.EXECUTORS.GetNames())
			os.Exit(1)
		},
	}

	names := actions.EXECUTORS.GetNames()
	for i := 0; i < len(names); i++ {
		executor := actions.EXECUTORS.GetByName(names[i])
		cmd := buildCommand(executor)
		rootCmd.AddCommand(cmd)
	}

	rootCmd.AddCommand(server)
	rootCmd.Execute()
}
