package repositories

import "batch/internal/domain/models"

type BatchFilesRepository interface {
	GetBatchFilesOfDay() ([]*models.BatchFile, error)
	InsertBatchFile(file *models.BatchFile) error
}
