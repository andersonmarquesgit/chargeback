package objectstorage

type Downloader interface {
	DownloadFile(localPath string, objectName string) error
}
