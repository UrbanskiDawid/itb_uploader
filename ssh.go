package main

import (
	"fmt"

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

func remote(user string, pass string, hostname string, port int) string {
	var cmd string = "ls /;"

	client, session, err := connectToHost(user, pass, fmt.Sprintf("%s:%d", hostname, port))
	if err != nil {
		panic(err)
	}
	out, err := session.CombinedOutput(cmd)
	if err != nil {
		panic(err)
	}
	client.Close()

	return string(out)
}
