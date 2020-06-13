package main

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"
)

func connectToHost(user, pass, host string) (*ssh.Client, *ssh.Session, error) {

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(pass)},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, err
	}

	return client, session, nil
}

func executeSSH(cmd string, serverName string) (string, error) {

	server := getServerByNickName(serverName)

	client, session, err := connectToHost(server.Auth.User, server.Auth.Pass, fmt.Sprintf("%s:%d", server.Host, server.Port))
	if err != nil {
		panic(err)
	}

	var out bytes.Buffer
	session.Stdin = strings.NewReader("")
	session.Stdout = &out
	err = session.Start(cmd)
	session.Wait()
	client.Close()
	return out.String(), err
}
