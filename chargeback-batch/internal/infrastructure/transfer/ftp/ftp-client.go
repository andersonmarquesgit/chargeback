package ftp

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"os"
	"time"
)

type FTPClient struct {
	conn *ftp.ServerConn
}

func NewFTPClient(host string, port int, username, password string) (*FTPClient, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := ftp.Dial(address, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to FTP server: %w", err)
	}

	if err := conn.Login(username, password); err != nil {
		return nil, fmt.Errorf("failed to login to FTP server: %w", err)
	}

	return &FTPClient{conn: conn}, nil
}

func (c *FTPClient) Upload(localPath, remotePath string) error {
	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file %s: %w", localPath, err)
	}
	defer file.Close()

	err = c.conn.Stor(remotePath, file)
	if err != nil {
		return fmt.Errorf("failed to upload file to FTP server: %w", err)
	}

	return nil
}

func (c *FTPClient) Close() error {
	if c.conn != nil {
		return c.conn.Quit()
	}
	return nil
}
