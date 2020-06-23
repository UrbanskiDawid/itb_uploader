package base

import "os"

type Action interface {
	Execute() (string, string, error)
	UploadFile(localFile string) (error, string)
	DownloadFile(localFile *os.File) (error, string)
	GetDescription() Description
}
