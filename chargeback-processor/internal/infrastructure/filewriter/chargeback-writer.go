package filewriter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"processor/internal/domain/models"
	"sync"
	"time"
)

type ChargebackWriter struct {
	mu            sync.Mutex
	currentFile   *os.File
	currentFileID string
	records       []models.Chargeback
	lastFlush     time.Time
	maxRecords    int
	maxDuration   time.Duration
	directory     string
}

func NewChargebackWriter(directory string, maxRecords int, maxDuration time.Duration) *ChargebackWriter {
	return &ChargebackWriter{
		records:     make([]models.Chargeback, 0),
		maxRecords:  maxRecords,
		maxDuration: maxDuration,
		directory:   directory,
		lastFlush:   time.Now(),
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
	if w.currentFile != nil {
		w.currentFile.Close()
	}

	// Garante que o diretório existe
	err := os.MkdirAll(w.directory, 0755)
	if err != nil {
		return fmt.Errorf("could not create output directory %s: %w", w.directory, err)
	}

	filename := fmt.Sprintf("cb_batch_%d.ndjson", time.Now().UnixNano())
	fullPath := filepath.Join(w.directory, filename)

	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("could not create file %s: %w", fullPath, err)
	}

	w.currentFile = file
	w.currentFileID = filename
	w.records = make([]models.Chargeback, 0)
	w.lastFlush = time.Now()
	return nil
}
