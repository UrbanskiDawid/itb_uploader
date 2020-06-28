package sshFiles

import (
	"errors"
	"fmt"

	"github.com/UrbanskiDawid/itb_uploader/logging"
	"github.com/hnakamur/go-scp"

	"golang.org/x/crypto/ssh"
)

func UploadFileScp(localFileName string, conn *ssh.Client, remoteFileName string) error {

	clientScp := scp.NewSCP(conn)
	if clientScp == nil {
		logging.LogConsole(fmt.Sprint("error: SCP fail to connect"))
		return errors.New("no SCP")
	}

	clientScp.SendFile(localFileName, remoteFileName)
	return nil

}

func DownloadFileScp(localFileName string, conn *ssh.Client, remoteFile string) error {
	clientScp := scp.NewSCP(conn)
	if clientScp == nil {
		logging.LogConsole(fmt.Sprint("error: SCP fail to connect"))
		return errors.New("no SCP")
	}

	err := clientScp.ReceiveFile(remoteFile, localFileName)
	if err != nil {
		logging.LogConsole(fmt.Sprintf("error: SCP fail while copying file", err))
		return err
	}
	return nil
}
