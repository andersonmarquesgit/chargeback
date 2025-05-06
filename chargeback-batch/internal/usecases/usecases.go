package usecases

// UseCases struct to store the providers use case and other use cases
type UseCases struct {
	ChargebackBatchEventUseCase *ChargebackBatchEventUseCase
}

// NewUseCases struct
func NewUseCases(chargebackBatchEventUseCase *ChargebackBatchEventUseCase) *UseCases {
	return &UseCases{
		ChargebackBatchEventUseCase: chargebackBatchEventUseCase,
	}
}
