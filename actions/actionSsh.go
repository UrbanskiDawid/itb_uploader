package actions

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/UrbanskiDawid/itb_uploader/actions/base"
	"github.com/UrbanskiDawid/itb_uploader/logging"
	scp "github.com/hnakamur/go-scp"
	"github.com/pkg/sftp"

	"golang.org/x/crypto/ssh"
)

type actionSsh struct {
	desc   base.Description
	server base.Server
	//todo ssh.ClientConfig
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func getAuthMethod(pass string) (ssh.AuthMethod, error) {

	if fileExists(pass) {
		key, err := ioutil.ReadFile(pass)
		if err != nil {
			return nil, err
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, err
		}

		return ssh.PublicKeys(signer), nil
	}

	if strings.HasPrefix(pass, "-----BEGIN RSA PRIVATE KEY-----") {
		signer, err := ssh.ParsePrivateKey([]byte(pass))
		if err != nil {
			return nil, err
		}
		return ssh.PublicKeys(signer), nil
	}

	return ssh.Password(pass), nil
}

func buildClientConfig(server base.Server) (*ssh.ClientConfig, error) {

	auth, err := getAuthMethod(server.Auth.Pass)
	if err != nil {
		logging.Log.Println("getAuthMethod failed")
		return nil, err
	}

	sshConfig := &ssh.ClientConfig{
		User:            server.Auth.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            []ssh.AuthMethod{auth},
	}

	return sshConfig, nil
}

func configureSSHforServer(server base.Server) (*ssh.Client, error) {

	sshConfig, err := buildClientConfig(server)
	if err != nil {
		logging.Log.Println("getAuthMethod failed")
		return nil, err
	}

	host := fmt.Sprintf("%s:%d", server.Host, server.Port)
	return ssh.Dial("tcp", host, sshConfig)
}

func (e actionSsh) Execute() (string, string, error) {

	cmd := e.desc.Cmd

	serverName := e.desc.Server

	logging.Log.Println("executeSSH", serverName, "cmd", cmd)

	//COMMON-------
	client, err := configureSSHforServer(e.server)
	if err != nil {
		logging.Log.Print("executeSSH configureSSHforServer fail")
		return "", "", err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		logging.Log.Print("executeSSH NewSession fail")
		return "", "", err
	}
	defer session.Close()
	//---

	var out bytes.Buffer
	var outErr bytes.Buffer
	session.Stdin = strings.NewReader("")
	session.Stdout = &out
	session.Stderr = &outErr

	err = session.Start(cmd)
	if err != nil {
		logging.Log.Print("executeSSH start fail")
		return out.String(), outErr.String(), err
	}
	defer session.Close()
	session.Wait()

	return out.String(), outErr.String(), err
}

func buildSftpClient(server base.Server) (*sftp.Client, error) {
	// connect
	conn, err := configureSSHforServer(server)
	if err != nil {
		logging.LogConsole("error buildSftpClient fail")
		return nil, err
	}
	defer conn.Close()

	// create new SFTP client
	client, err := sftp.NewClient(conn)
	if err != nil {
		logging.LogConsole("error SftpNewClient")
		return nil, err
	}
	return client, nil
}

func buildScpClient(server base.Server) (*scp.SCP, error) {

	client, err := configureSSHforServer(server)
	if err != nil {
		logging.LogConsole("error buildSftpClient fail")
		return nil, err
	}

	scpClient := scp.NewSCP(client)
	if scpClient == nil {
		logging.LogConsole("error ScpNewClient")
		return nil, err
	}

	return scpClient, err
}

func uploadFileSftp(localFileName string, server base.Server, remoteFileName string) error {

	//COMMON-------
	conn, err := configureSSHforServer(server)
	if err != nil {
		logging.Log.Print("sendFile configureSSHforServer fail")
		return err
	}
	defer conn.Close()
	//---

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

func uploadFileScp(localFileName string, server base.Server, remoteFileName string) error {

	clientScp, err := buildScpClient(server)
	if err != nil {
		logging.LogConsole(fmt.Sprintf("error: SCP fail to connect", err))
		return err
	}
	clientScp.SendFile(localFileName, remoteFileName)
	return nil

}

//UploadFile send local file
func (e actionSsh) UploadFile(localFileName string) error {

	remoteFileName := e.desc.FileTarget

	err := uploadFileSftp(localFileName, e.server, remoteFileName)
	if err == nil {
		return nil
	}
	//SFTP CLIENT
	return uploadFileScp(localFileName, e.server, remoteFileName)
}

//DownloadFile get remote file
func downloadFileSftp(localFileName string, server base.Server, remoteFile string) error {

	// connect
	conn, err := configureSSHforServer(server)
	if err != nil {
		logging.LogConsole(fmt.Sprint("downloadFileSftp configureSSHforServer fail", err))
		return err
	}
	defer conn.Close()

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

func downloadFileScp(localFileName string, server base.Server, remoteFile string) error {
	clientScp, err := buildScpClient(server)
	if err != nil {
		logging.LogConsole(fmt.Sprintf("error: SCP fail to connect", err))
		return err
	}

	err = clientScp.ReceiveFile(localFileName, remoteFile)
	if err != nil {
		logging.LogConsole(fmt.Sprintf("error: SCP fail while copying file", err))
		return err
	}
	return nil
}

func (e actionSsh) DownloadFile(localFile string) error {

	remoteFile := e.desc.FileDownload
	err := downloadFileSftp(localFile, e.server, remoteFile)
	if err == nil {
		return nil
	}
	return downloadFileScp(localFile, e.server, remoteFile)
}

func (e actionSsh) GetDescription() base.Description {
	return e.desc
}
