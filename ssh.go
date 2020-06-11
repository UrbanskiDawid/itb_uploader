package main

import (
	"fmt"
	"os"

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

type RemoteCmd struct {
	running bool
	out     string
}

var remoteCmd RemoteCmd

func remote(cmd string) {

	config := loadConfig()

	// wait for the program to end in a goroutine
	go func() {
		client, session, err := connectToHost(config.user, config.pass, config.host)
		if err != nil {
			panic(err)
		}
		remoteCmd.running = true
		out, err := session.CombinedOutput(cmd)
		if err != nil {
			panic(err)
		}
		client.Close()
		remoteCmd.running = false
		remoteCmd.out = string(out)
	}()
}
