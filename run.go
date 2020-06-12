package main

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)

type MyCmd struct {
	cmd *exec.Cmd
	out bytes.Buffer
}

func run(cmd string) MyCmd {

	var myCmd MyCmd

	myCmd.cmd = exec.Command(cmd)
	myCmd.cmd.Stdin = strings.NewReader("")
	myCmd.cmd.Stdout = &myCmd.out

	err := myCmd.cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	return myCmd
}
