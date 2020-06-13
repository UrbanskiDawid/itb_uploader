package workers

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)

func executeLocal(cmd string) (string, string, error) {

	log.Println("executeLocal", cmd)

	var out bytes.Buffer
	var outErr bytes.Buffer

	s := strings.Split(cmd, " ")

	app := s[0]
	args := s[1:]

	exe := exec.Command(app, args...)
	exe.Stdin = strings.NewReader("")
	exe.Stdout = &out
	exe.Stderr = &outErr

	err := exe.Start()
	if err != nil {
		log.Println("executeLocal failed to start", err)
		return "", "", err
	}

	err = exe.Wait()
	if err != nil {
		print(outErr.String())
	}

	log.Printf("executeLocal failed at wait %s\n", err)
	return out.String(), outErr.String(), err
}
