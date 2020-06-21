package base

type Action interface {
	Execute() (string, string, error)
	UploadFile(localFile string) (error, string)
	DownloadFile(localFile string) (error, string)
	GetDescription() Description
}
