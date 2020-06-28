package sshFiles

import (
	"fmt"
	"io"
	"os"

	"github.com/UrbanskiDawid/itb_uploader/logging"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

//DownloadFile get remote file
func DownloadFileSftp(localFileName string, conn *ssh.Client, remoteFile string) error {

	// create new SFTP client
	client, err := sftp.NewClient(conn)
	if err != nil {
		logging.LogConsole(fmt.Sprint("downloadFileSftp fail to create client", err))
		return err
	}
	defer client.Close()

	// create destination file
	dstFile, err := os.Create(localFileName)
	if err != nil {
		logging.LogConsole(fmt.Sprint("downloadFileSftp fail to create localFileName", err))
		return err
	}
	defer dstFile.Close()

	// open source file
	srcFile, err := client.Open(remoteFile)
	if err != nil {
		logging.LogConsole(fmt.Sprint("downloadFileSftp fail to open remotefile", err))
		return err
	}

	// copy source file to destination file
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		logging.LogConsole(fmt.Sprint("downloadFileSftp fail to copy files", err))
		return err
	}
	//fmt.Printf("%d bytes copied\n", bytes)

	// flush in-memory copy
	err = dstFile.Sync()
	if err != nil {
		logging.LogConsole(fmt.Sprint("downloadFileSftp fail sync", err))
		return err
	}

	return nil
}

func UploadFileSftp(localFileName string, conn *ssh.Client, remoteFileName string) error {

	// create new SFTP client
	client, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}
	defer client.Close()

	// create destination file
	dstFile, err := client.Create(remoteFileName)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// create source file
	srcFile, err := os.Open(localFileName)
	if err != nil {
		return err
	}

	// copy source file to destination file
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	//fmt.Printf("%d bytes copied\n", bytes)
	return nil
}
