package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

type SSHconfig struct {
	user string
	pass string
	host string
	port int
}

func loadConfig() SSHconfig {

	var ret SSHconfig
	ret.user = os.Getenv("SSH_USER")
	ret.pass = os.Getenv("SSH_PASS")
	ret.host = fmt.Sprintf("%s:%s", os.Getenv("SSH_HOST"), os.Getenv("SSH_PORT"))

	if ret.user == "" || ret.pass == "" || ret.host == "" {
		panic("enviroment configuration missing")
	}

	return ret
}

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

func Remote(cmd string) (string, error) {

	config := loadConfig()

	client, session, err := connectToHost(config.user, config.pass, config.host)
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
