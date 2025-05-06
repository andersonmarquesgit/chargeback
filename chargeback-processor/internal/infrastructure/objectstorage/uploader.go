package objectstorage

type Uploader interface {
	UploadFile(localPath string, objectName string) error
}
