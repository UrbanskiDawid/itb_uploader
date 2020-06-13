package main

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)

func executeLocal(cmd string) (string, error) {

	var out bytes.Buffer

	s := strings.Split(cmd, " ")

	app := s[0]
	args := s[1:]

	exe := exec.Command(app, args...)
	exe.Stdin = strings.NewReader("")
	exe.Stdout = &out

	err := exe.Start()
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	err = exe.Wait()
	if err != nil {
		return "", err
	}

	return out.String(), err
}
