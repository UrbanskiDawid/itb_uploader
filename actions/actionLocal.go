package actions

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/UrbanskiDawid/itb_uploader/actions/base"

	"github.com/UrbanskiDawid/itb_uploader/logging"
)

type ActionLocal struct {
	desc base.Description
}

func (e ActionLocal) Execute() (string, string, error) {

	cmd := e.desc.Cmd

	logging.LogConsole("executeLocal: " + cmd)

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
		logging.Logger.Print("executeLocal failed to start", err)
		return "", "", err
	}

	err = exe.Wait()
	if err != nil {
		logging.Logger.Printf("executeLocal failed at wait %s\n", outErr.String())
	}

	return out.String(), outErr.String(), err
}

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func copyFileLocal(src, dst string) error {

	if !fileExists(src) {
		return errors.New("file '" + src + "' does not exist")
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func (e ActionLocal) UploadFile(localFileName string) error {
	dst := e.desc.FileTarget

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	fn, err := os.Open(localFileName)
	if err != nil {
		return err
	}
	defer fn.Close()

	_, err = io.Copy(out, fn)
	if err != nil {
		return err
	}

	return nil
}

func (e ActionLocal) DownloadFile(outFileName string) error {
	src := e.desc.FileDownload

	if !fileExists(src) {
		return errors.New("file '" + src + "' does not exist")
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	fn, err := os.OpenFile(outFileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer fn.Close()

	_, err = io.Copy(fn, in)
	if err != nil {
		return err
	}

	return nil
}

func (e ActionLocal) GetDescription() base.Description {
	return e.desc
}
