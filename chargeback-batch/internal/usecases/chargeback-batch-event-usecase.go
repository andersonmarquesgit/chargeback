package usecases

import (
	"batch/internal/domain/models"
	"batch/internal/domain/repositories"
	"batch/internal/infrastructure/logging"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ChargebackBatchEventUseCase struct {
	BatchFilesRepository repositories.BatchFilesRepository
}

func NewChargebackBatchEventUseCase(batchFilesRepository repositories.BatchFilesRepository) *ChargebackBatchEventUseCase {
	return &ChargebackBatchEventUseCase{
		BatchFilesRepository: batchFilesRepository,
	}
}

func (uc *ChargebackBatchEventUseCase) Process(msg amqp.Delivery) error {
	var bf models.BatchFile
	if err := json.Unmarshal(msg.Body, &bf); err != nil {
		logging.Infof("failed to unmarshal batch file: %v", err)
		return err
	}

	// Persistir os detalhes do batch file, url no object storage, etc.
	err := uc.BatchFilesRepository.InsertBatchFile(&bf)
	if err != nil {
		logging.Infof("failed to persist batch file: %v", err)
		return err
	}

	return nil
}
