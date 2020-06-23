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

func (e ActionLocal) UploadFile(localFile io.Reader) (error, string) {
	dst := e.desc.FileTarget

	out, err := os.Create(dst)
	if err != nil {
		return err, ""
	}
	defer out.Close()

	_, err = io.Copy(out, localFile)
	if err != nil {
		return err, ""
	}

	return nil, dst
}

func (e ActionLocal) DownloadFile(outFile *os.File) (error, string) {
	src := e.desc.FileDownload

	defer outFile.Close()

	if !fileExists(src) {
		return errors.New("file '" + src + "' does not exist"), ""
	}

	in, err := os.Open(src)
	if err != nil {
		return err, ""
	}
	defer in.Close()

	_, err = io.Copy(outFile, in)
	if err != nil {
		return err, ""
	}

	return nil, src
}

func (e ActionLocal) GetDescription() base.Description {
	return e.desc
}
