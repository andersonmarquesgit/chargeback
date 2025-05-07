package ftp

type Client interface {
	Upload(localPath, remotePath string) error
	Close() error
}
