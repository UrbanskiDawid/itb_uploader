package base

import (
	"io"
	"os"
)

type Action interface {
	Execute() (string, string, error)
	UploadFile(io.Reader) (error, string)
	DownloadFile(localFile *os.File) (error, string)
	GetDescription() Description
}
