package usecases

// UseCases struct to store the providers use case and other use cases
type UseCases struct {
	ChargebackOpenedEventUseCase *ChargebackOpenedEventUseCase
}

// NewUseCases struct
func NewUseCases(chargebackOpenedUseCase *ChargebackOpenedEventUseCase) *UseCases {
	return &UseCases{
		ChargebackOpenedEventUseCase: chargebackOpenedUseCase,
	}
}
