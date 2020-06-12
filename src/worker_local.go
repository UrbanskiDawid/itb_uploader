package main

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)

func Run(cmd string) (string, error) {

	var out bytes.Buffer

	exe := exec.Command(cmd)
	exe.Stdin = strings.NewReader("")
	exe.Stdout = &out

	err := exe.Start()
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	err = exe.Wait()
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return out.String(), err
}
