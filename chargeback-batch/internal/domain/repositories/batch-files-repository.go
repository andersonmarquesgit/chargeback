package repositories

import "batch/internal/domain/models"

type BatchFilesRepository interface {
	GetBatchFilesOfDay(batchMaxFiles int) ([]*models.BatchFile, error)
	InsertBatchFile(file *models.BatchFile) error
	MarkAsFailed(fileID string) error
	MarkAsSent(fileID string) error
}
