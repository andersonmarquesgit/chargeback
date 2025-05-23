package postgres

import (
	"batch/internal/domain/models"
	"batch/internal/domain/repositories"
	"database/sql"
	"fmt"
)

type BatchFilesRepositoryPostgres struct {
	DB *sql.DB
}

func NewBatchFilesRepositoryPostgres(db *sql.DB) repositories.BatchFilesRepository {
	return &BatchFilesRepositoryPostgres{
		DB: db,
	}
}

func (r *BatchFilesRepositoryPostgres) GetBatchFilesOfDay(batchMaxFiles int) ([]*models.BatchFile, error) {
	query := `
		SELECT file_id, file_url, created_at, record_count, status, sent_at, retry_count, last_attempt_at
		FROM batch_files
		WHERE 
			(status = 'ready' OR status = 'failed')
			AND retry_count <= 3
		ORDER BY created_at ASC
		LIMIT $1;
	`

	rows, err := r.DB.Query(query, batchMaxFiles)
	if err != nil {
		return nil, fmt.Errorf("could not query batch_files: %w", err)
	}
	defer rows.Close()

	var batchFiles []*models.BatchFile
	for rows.Next() {
		var bf models.BatchFile
		err := rows.Scan(
			&bf.FileID,
			&bf.FileURL,
			&bf.CreatedAt,
			&bf.RecordCount,
			&bf.Status,
			&bf.SentAt,
			&bf.RetryCount,
			&bf.LastAttemptAt,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan row: %w", err)
		}
		batchFiles = append(batchFiles, &bf)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return batchFiles, nil
}

func (r *BatchFilesRepositoryPostgres) InsertBatchFile(file *models.BatchFile) error {
	query := `
		INSERT INTO batch_files (
			file_id,
			file_url,
			created_at,
			record_count,
			status,
			sent_at,
			retry_count,
			last_attempt_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.DB.Exec(
		query,
		file.FileID,
		file.FileURL,
		file.CreatedAt,
		file.RecordCount,
		file.Status,
		file.SentAt,
		file.RetryCount,
		file.LastAttemptAt,
	)
	if err != nil {
		return fmt.Errorf("could not insert batch file: %w", err)
	}

	return nil
}

func (r *BatchFilesRepositoryPostgres) MarkAsSent(fileID string) error {
	query := `
		UPDATE batch_files
		SET status = 'sent',
			sent_at = NOW(),
			last_attempt_at = NOW()
		WHERE file_id = $1
	`
	_, err := r.DB.Exec(query, fileID)
	if err != nil {
		return fmt.Errorf("could not mark file as sent: %w", err)
	}
	return nil
}

func (r *BatchFilesRepositoryPostgres) MarkAsFailed(fileID string) error {
	query := `
		UPDATE batch_files
		SET status = 'failed',
			retry_count = retry_count + 1,
			last_attempt_at = NOW()
		WHERE file_id = $1
	`
	_, err := r.DB.Exec(query, fileID)
	if err != nil {
		return fmt.Errorf("could not mark file as failed: %w", err)
	}
	return nil
}
