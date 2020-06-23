package base

type Action interface {
	Execute() (string, string, error)
	UploadFile(string) error
	DownloadFile(string) error
	GetDescription() Description
}
