package scheduler

import (
	"batch/internal/domain/repositories"
	"batch/internal/infrastructure/logging"
	"batch/internal/infrastructure/objectstorage"
	"batch/internal/infrastructure/transfer/ftp"
	"path/filepath"
	"time"
)

func StartScheduler(repo repositories.BatchFilesRepository, downloader objectstorage.Downloader, ftpClient ftp.Client, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			processBatchFiles(repo, downloader, ftpClient)
		}
	}
}

func processBatchFiles(repo repositories.BatchFilesRepository, downloader objectstorage.Downloader, ftpClient ftp.Client) {
	files, err := repo.GetBatchFilesOfDay()
	if err != nil {
		logging.Infof("Failed to fetch batch files: %v", err)
		return
	}

	for _, f := range files {
		localPath := filepath.Join("/tmp/chargebacks", f.FileID)
		err := downloader.DownloadFile(localPath, f.FileID)
		if err != nil {
			logging.Infof("Failed to download file %s: %v", f.FileID, err)
			repo.MarkAsFailed(f.FileID)
			continue
		}

		err = ftpClient.Upload(localPath, f.FileID)
		if err != nil {
			logging.Infof("Failed to send file %s via FTP: %v", f.FileID, err)
			repo.MarkAsFailed(f.FileID)
			continue
		}

		repo.MarkAsSent(f.FileID)
		logging.Infof("Successfully sent file %s via FTP", f.FileID)
	}
}
