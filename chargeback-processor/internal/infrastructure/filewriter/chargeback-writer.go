package filewriter

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"processor/internal/domain/models"
	"processor/internal/infrastructure/logging"
	"processor/internal/infrastructure/objectstorage"
	"processor/internal/infrastructure/rabbitmq/producers"
	"sync"
	"time"
)

type ChargebackWriter struct {
	mu                 sync.Mutex
	currentFile        *os.File
	currentFileID      string
	records            []models.Chargeback
	lastFlush          time.Time
	maxRecords         int
	maxDuration        time.Duration
	directory          string
	uploader           objectstorage.Uploader
	batchEventProducer *producers.Producer
}

func NewChargebackWriter(directory string, maxRecords int, maxDuration time.Duration, uploader objectstorage.Uploader, producer *producers.Producer) *ChargebackWriter {
	return &ChargebackWriter{
		records:            make([]models.Chargeback, 0),
		maxRecords:         maxRecords,
		maxDuration:        maxDuration,
		directory:          directory,
		lastFlush:          time.Now(),
		uploader:           uploader,
		batchEventProducer: producer,
	}
}

func (w *ChargebackWriter) Write(cb models.Chargeback) (string, error) {
	// Garante que temos um arquivo aberto
	if w.currentFile == nil {
		if err := w.rotateFile(); err != nil {
			return "", fmt.Errorf("failed to open new file: %w", err)
		}
	}

	// Adiciona o chargeback à memória
	w.records = append(w.records, cb)

	// Escreve a linha no arquivo atual
	line, err := json.Marshal(cb)
	if err != nil {
		return "", fmt.Errorf("failed to marshal chargeback: %w", err)
	}
	if _, err := w.currentFile.Write(append(line, '\n')); err != nil {
		return "", fmt.Errorf("failed to write chargeback to file: %w", err)
	}

	// Verifica se deve rotacionar
	if len(w.records) >= w.maxRecords || time.Since(w.lastFlush) >= w.maxDuration {
		if err := w.rotateFile(); err != nil {
			return "", fmt.Errorf("failed to rotate file: %w", err)
		}
	}

	return w.currentFileID, nil
}

func (w *ChargebackWriter) rotateFile() error {
	var prevFileID string
	var recordCount int

	// Se existir um arquivo atual, faz upload para o MinIO antes de criar o novo
	if w.currentFile != nil {
		prevFileID = w.currentFileID
		recordCount = len(w.records)

		w.currentFile.Close()

		fullPath := filepath.Join(w.directory, w.currentFileID)
		if err := w.uploader.UploadFile(fullPath, w.currentFileID); err != nil {
			return fmt.Errorf("failed to upload chargeback file to object storage: %w", err)
		}

		if err := w.notifyBatchReady(w.currentFileID, len(w.records)); err != nil {
			return fmt.Errorf("failed to publish batch ready event: %w", err)
		}

		logging.Infof("Rotating file with %d records: %s", recordCount, prevFileID)
	}

	// Garante que o diretório existe
	if err := os.MkdirAll(w.directory, 0755); err != nil {
		return fmt.Errorf("could not create output directory %s: %w", w.directory, err)
	}

	// Cria um novo arquivo com timestamp único
	filename := fmt.Sprintf("cb_batch_%d.ndjson", time.Now().UnixNano())
	fullPath := filepath.Join(w.directory, filename)

	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("could not create file %s: %w", fullPath, err)
	}

	// Atualiza estado interno do writer
	w.currentFile = file
	w.currentFileID = filename
	w.records = make([]models.Chargeback, 0)
	w.lastFlush = time.Now()

	return nil
}

func (w *ChargebackWriter) notifyBatchReady(fileID string, recordCount int) error {
	event := map[string]interface{}{
		"file_id":      fileID,
		"file_url":     fmt.Sprintf("http://minio:9000/%s/%s", w.uploader.GetBucketName(), fileID),
		"created_at":   time.Now().UTC().Format(time.RFC3339),
		"record_count": recordCount,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("could not marshal batch ready event: %w", err)
	}

	traceID := uuid.NewString()
	return w.batchEventProducer.PublishChargebackBatchEvent(payload, traceID)
}

func (w *ChargebackWriter) MaybeFlush() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if len(w.records) == 0 || w.currentFile == nil {
		return
	}

	if time.Since(w.lastFlush) >= w.maxDuration {
		_ = w.rotateFile()
	}
}
