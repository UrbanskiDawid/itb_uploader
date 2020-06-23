package main

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/UrbanskiDawid/itb_uploader/actions"
	"github.com/UrbanskiDawid/itb_uploader/actions/base"
	"github.com/UrbanskiDawid/itb_uploader/views"
	"github.com/spf13/cobra"
)

func buildCommand(action base.Action) *cobra.Command {

	description := action.GetDescription()

	var cmd = &cobra.Command{
		Use: action.GetDescription().Name,
		Run: func(cmd *cobra.Command, args []string) {

			if description.HasUploadFile() {

				localFileName := args[0]

				err := action.UploadFile(localFileName)
				if err != nil {
					print(err)
					os.Exit(1)
				}
				fmt.Printf("file %s sent to %s@%s", localFileName, action.GetDescription().Server, action.GetDescription().FileTarget)
			}

			if description.HasCommand() {
				stdOut, stdErr, err := action.Execute()
				print(stdOut)
				if err != nil {
					print(stdErr)
					os.Exit(1)
				}
			}

			if description.HasDownloadFile() {

				remoteFileBaseName := path.Base(action.GetDescription().FileDownload)
				outFileName := path.Join(".", remoteFileBaseName)

				err := action.DownloadFile(outFileName)
				if err != nil {
					print("ERROR", err)
					os.Exit(1)
				}
				println("new file:", outFileName)
			}
		},
	}

	short := description.Name
	if description.HasUploadFile() {
		cmd.Args = cobra.ExactArgs(1)
		short = fmt.Sprintf("%s [file]", short)
	}
	if description.HasDownloadFile() {
		short = fmt.Sprintf("%s this cmd will download file", short)
	}
	if description.HasCommand() {
		short = fmt.Sprintf("%s this cmd will run command '%s'", short, description.Cmd)
	}
	cmd.Short = short
	return cmd
}

func runCli() {

	var rootCmd = &cobra.Command{
		//	Use: "app"
	}

	var server = &cobra.Command{
		Use:   "server [start server]",
		Short: "server port",
		Long:  `run this app in server mode.`,
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				port, _ = strconv.ParseUint(args[0], 10, 64)
			}
			err := views.StartServer(port)
			if err != nil {
				fmt.Println("server failed", err)
				logger.Fatal("server failed", err)
				os.Exit(1)
			}
			os.Exit(0)
		},
	}

	names := actions.ACTIONS.GetNames()
	for i := 0; i < len(names); i++ {
		executor := actions.ACTIONS.GetByName(names[i])
		cmd := buildCommand(executor)
		rootCmd.AddCommand(cmd)
	}

	rootCmd.AddCommand(server)
	rootCmd.Execute()
}
