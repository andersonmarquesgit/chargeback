package usecases

// UseCases struct to store the customer use case and other use cases
type UseCases struct {
	ChargebackOpenedUseCase *ChargebackOpenedUseCase
}

// NewUseCases struct
func NewUseCases(chargebackOpenedUseCase *ChargebackOpenedUseCase) *UseCases {
	return &UseCases{
		ChargebackOpenedUseCase: chargebackOpenedUseCase,
	}
}
