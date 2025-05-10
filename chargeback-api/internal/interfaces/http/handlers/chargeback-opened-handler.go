package handlers

import (
	"api/internal/interfaces/dto"
	"api/internal/presentation"
	"api/internal/usecases"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ChargebackOpenedHandler struct {
	chargebackOpenedUseCase *usecases.ChargebackOpenedUseCase
}

func NewChargebackOpenedHandler(chargebackOpenedUseCase *usecases.ChargebackOpenedUseCase) *ChargebackOpenedHandler {
	return &ChargebackOpenedHandler{
		chargebackOpenedUseCase: chargebackOpenedUseCase,
	}
}

// CreateChargeback godoc
// @Summary Open a chargeback for user and transaction
// @Description Verify if the chargeback exists using idempotency with user ID and transaction ID.<br><br>- **If it does not exist:** Sends a message to the queue (`chargeback-opened`) to create a new chargeback in the processor. <br><br>- **Returns:** <br><br>`202 Accepted` <br><br>`"message": "Chargeback sent to processor successfully"` <br><br><br>- **If it already exists:** <br><br>- **Returns:** <br><br>`200 OK` <br><br>`"message": "Chargeback already exists"`
// @Tags chargeback
// @Accept json
// @Produce json
// @Param chargeback body dto.ChargebackRequest true "Data of the chargeback"
// @Success 202 {object} presentation.JSONResponse
// @Success 200 {object} presentation.JSONResponse
// @Failure 400 {object} presentation.JSONResponse
// @Failure 500 {object} presentation.JSONResponse
// @Router /v1/chargebacks [post]
func (h *ChargebackOpenedHandler) CreateChargeback(w http.ResponseWriter, r *http.Request) {
	var req dto.ChargebackRequest

	// Decode the request body into the chargeback request struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		presentation.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Validate the chargeback request struct
	if err := validate.Struct(req); err != nil {
		presentation.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Get the request-trace-id from context
	traceID := r.Context().Value("request-trace-id").(string)

	// Create chargeback
	chargeback, err := h.chargebackOpenedUseCase.CreateChargeback(req, traceID)
	if err != nil {
		presentation.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	status := http.StatusAccepted
	message := "Chargeback sent to processor successfully"
	if chargeback.Exists {
		status = http.StatusOK
		message = "Chargeback already exists"
	}

	// Return the chargeback send to processor
	presentation.WriteJSON(w, status, presentation.JSONResponse{
		Error:   false,
		Message: message,
		Data:    nil,
	})
}
