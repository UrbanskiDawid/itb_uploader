package sshBackend

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/UrbanskiDawid/itb_uploader/actions/base"
	"github.com/UrbanskiDawid/itb_uploader/logging"

	"golang.org/x/crypto/ssh"
)

type ActionSsh struct {
	desc   base.Description
	server base.Server
	config ssh.ClientConfig
}

func BuildActionSsh(description base.Description, server base.Server) ActionSsh {
	client, err := buildClientConfig(server)
	if err != nil {
		panic("server " + server.NickName + " configuration error")
	}
	return ActionSsh{description, server, *client}
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
		logging.Logger.Println("getAuthMethod failed")
		return nil, err
	}

	sshConfig := &ssh.ClientConfig{
		User:            server.Auth.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            []ssh.AuthMethod{auth},
	}

	return sshConfig, nil
}

func configureSSHforServer(action ActionSsh) (*ssh.Client, error) {
	host := fmt.Sprintf("%s:%d", action.server.Host, action.server.Port)
	return ssh.Dial("tcp", host, &action.config)
}

func (e ActionSsh) Execute() (string, string, error) {

	cmd := e.desc.Cmd

	serverName := e.desc.Server

	logging.Logger.Println("executeSSH", serverName, "cmd", cmd)

	//COMMON-------
	client, err := configureSSHforServer(e)
	if err != nil {
		logging.Logger.Print("executeSSH configureSSHforServer fail")
		return "", "", err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		logging.Logger.Print("executeSSH NewSession fail")
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
		logging.Logger.Print("executeSSH start fail")
		return out.String(), outErr.String(), err
	}
	defer session.Close()
	session.Wait()

	return out.String(), outErr.String(), err
}

//UploadFile send local file
func (e ActionSsh) UploadFile(localFileName string) error {

	remoteFileName := e.desc.FileTarget

	client, err := configureSSHforServer(e)
	if err != nil {
		return err
	}

	err = UploadFileSftp(localFileName, client, remoteFileName)
	if err == nil {
		return nil
	}
	//SFTP CLIENT
	return UploadFileScp(localFileName, client, remoteFileName)
}

func (e ActionSsh) DownloadFile(localFile string) error {

	remoteFile := e.desc.FileDownload

	client, err := configureSSHforServer(e)
	if err != nil {
		return err
	}

	err = DownloadFileSftp(localFile, client, remoteFile)
	if err == nil {
		return nil
	}
	return DownloadFileScp(localFile, client, remoteFile)
}

func (e ActionSsh) GetDescription() base.Description {
	return e.desc
}
