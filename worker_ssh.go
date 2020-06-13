package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

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
			return nil, err //log.Fatalf("unable to read private key: %v", err)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, err //log.Fatalf("unable to parse private key: %v", err)
		}

		return ssh.PublicKeys(signer), nil

	}
	return ssh.Password(pass), nil
}

func configureSSHforServer(serverName string) (*ssh.Client, error) {

	server := getServerByNickName(serverName)

	auth, err := getAuthMethod(server.Auth.Pass)
	if err != nil {
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

func executeSSH(cmd string, serverName string) (string, error) {

	client, err := configureSSHforServer(serverName)
	if err != nil {
		println("executeSSH configureSSHforServer fail")
		return "", err
	}

	session, err := client.NewSession()
	if err != nil {
		println("executeSSH NewSession fail")
		return "", err
	}
	defer client.Close()

	var out bytes.Buffer
	session.Stdin = strings.NewReader("")
	session.Stdout = &out
	err = session.Start(cmd)
	if err != nil {
		println("executeSSH start fail")
		return "", err
	}
	defer session.Close()
	session.Wait()

	return out.String(), err
}
