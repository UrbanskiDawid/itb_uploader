package workers

import (
	"bytes"
	"os/exec"
	"strings"
	"logging"
)

func executeLocal(cmd string) (string, string, error) {

	logging.Log.Println("executeLocal", cmd)

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
		logging.Log.Println("executeLocal failed to start", err)
		return "", "", err
	}

	err = exe.Wait()
	if err != nil {
		logging.Log.Printf("executeLocal failed at wait %s\n", outErr.String())
	}

	return out.String(), outErr.String(), err
}
