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

func configureSSHforServer(server base.Server) (*ssh.Client, error) {

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

//UploadFile send local file
func (e actionSsh) UploadFile(localFile string) (error, string) {

	remoteFile := e.desc.FileTarget

	//COMMON-------
	conn, err := configureSSHforServer(e.server)
	if err != nil {
		logging.Log.Print("sendFile configureSSHforServer fail")
		return err, ""
	}
	defer conn.Close()
	//---

	// create new SFTP client
	client, err := sftp.NewClient(conn)
	if err != nil {
		return err, ""
	}
	defer client.Close()

	// create destination file
	dstFile, err := client.Create(remoteFile)
	if err != nil {
		return err, ""
	}
	defer dstFile.Close()

	// create source file
	srcFile, err := os.Open(localFile)
	if err != nil {
		return err, ""
	}

	// copy source file to destination file
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err, ""
	}
	//fmt.Printf("%d bytes copied\n", bytes)
	return nil, remoteFile
}

//DownloadFile get remote file
func (e actionSsh) DownloadFile(localFile string) (error, string) {

	remoteFile := e.desc.FileDownload

	// connect
	conn, err := configureSSHforServer(e.server)
	if err != nil {
		logging.Log.Print("sendFile configureSSHforServer fail")
		return err, ""
	}
	defer conn.Close()

	// create new SFTP client
	client, err := sftp.NewClient(conn)
	if err != nil {
		return err, ""
	}
	defer client.Close()

	// create destination file
	dstFile, err := os.Create(localFile)
	if err != nil {
		return err, ""
	}
	defer dstFile.Close()

	// open source file
	srcFile, err := client.Open(remoteFile)
	if err != nil {
		return err, ""
	}

	// copy source file to destination file
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err, ""
	}
	//fmt.Printf("%d bytes copied\n", bytes)

	// flush in-memory copy
	err = dstFile.Sync()
	if err != nil {
		return err, ""
	}

	return nil, remoteFile
}

func (e actionSsh) GetDescription() base.Description {
	return e.desc
}
