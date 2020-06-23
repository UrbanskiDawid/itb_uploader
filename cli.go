package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/UrbanskiDawid/itb_uploader/actions"
	"github.com/UrbanskiDawid/itb_uploader/actions/base"
	"github.com/UrbanskiDawid/itb_uploader/views"
	"github.com/spf13/cobra"
)

var tempFolder = "TEMP"

func openTmpFile() (*os.File, error) {
	//temp file
	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	tempPath := filepath.Join(path, tempFolder)
	os.MkdirAll(tempPath, os.ModePerm)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile(tempPath, "download-*")
	if err != nil {
		return nil, err
	}

	//defer os.Remove(tempFile.Name())
	return tempFile, nil
}

func buildCommand(action base.Action) *cobra.Command {

	description := action.GetDescription()

	var cmd = &cobra.Command{
		Use: action.GetDescription().Name,
		Run: func(cmd *cobra.Command, args []string) {

			if description.HasUploadFile() {
				err, _ := action.UploadFile(args[1])
				if err != nil {
					print(err)
					os.Exit(1)
				}
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

				tmpFile, err := openTmpFile()
				if err != nil {
					print("file download err", err)
				}
				defer tmpFile.Close()

				err, fn := action.DownloadFile(tmpFile)
				if err != nil {
					print("ERROR", err)
					os.Exit(1)
				}
				println("new file:", tmpFile.Name(), fn)
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
