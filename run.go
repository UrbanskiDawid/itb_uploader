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

var myCmd MyCmd

func run(cmd string) {

	myCmd.cmd = exec.Command(cmd)
	myCmd.cmd.Stdin = strings.NewReader("")
	myCmd.cmd.Stdout = &myCmd.out

	err := myCmd.cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

}
