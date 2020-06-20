package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/UrbanskiDawid/itb_uploader/actions"
	"github.com/UrbanskiDawid/itb_uploader/logging"
	"github.com/UrbanskiDawid/itb_uploader/views"
	"github.com/spf13/cobra"
)

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
			views.StartServer(port)
			os.Exit(1)
		},
	}

	names := actions.GetActionNames()
	for i := 0; i < len(names); i++ {

		var actionName = names[i]

		var cmd = &cobra.Command{
			Use: actionName,
			Run: func(cmd *cobra.Command, args []string) {

				if actions.IsActionWithUploadFile(actionName) {
					err := actions.UploadFile(actionName, args[1])
					if err != nil {
						print(err)
						os.Exit(1)
					}
				}

				stdOut, stdErr, err := actions.ExecuteAction(actionName)
				print(stdOut)
				if err != nil {
					print(stdErr)
					os.Exit(1)
				}

				if actions.IsActionWithDownloadFile(actionName) {

					targetFileName := "download"
					logging.Log.Printf("download %s to %s", actions.GetDownloadFileNameForAction(actionName), targetFileName)

					err := actions.DownloadFile(actionName, targetFileName)
					if err != nil {
						print("ERROR", err)
						print(stdErr)
						os.Exit(1)
					}
					println("new file:", targetFileName)
				}
			},
		}

		short := ""
		if actions.IsActionWithUploadFile(actionName) {
			cmd.Args = cobra.ExactArgs(1)
			short = fmt.Sprintf("%s [file]", actionName)
		}
		if actions.IsActionWithDownloadFile(actionName) {
			short = fmt.Sprintf("%sthis cmd will download file", short)
		}
		cmd.Short = short

		rootCmd.AddCommand(cmd)

	}

	rootCmd.AddCommand(server)
	rootCmd.Execute()
}
