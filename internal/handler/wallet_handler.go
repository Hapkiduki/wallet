package handler

import (
	"log/slog"
	"strings"
	"wallet/internal/usecase"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v3"
)

type WalletHandler struct {
	walletUsecase usecase.WalletUsecase
	logger        *slog.Logger
}

func NewWalletHandler(wu usecase.WalletUsecase, logger *slog.Logger) *WalletHandler {
	return &WalletHandler{walletUsecase: wu, logger: logger}
}

type RechargeRequest struct {
	WalletID string  `json:"wallet_id"`
	Amount   float64 `json:"amount"`
}

// @Summary Recharge a wallet
// @Description Adds a specified amount to a wallet's balance.
// @Tags wallets
// @Accept json
// @Produce json
// @Param wallet body RechargeRequest true "Recharge details"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /wallets/recharge [post]
func (h *WalletHandler) Recharge(c fiber.Ctx) error {
	var req RechargeRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse request"})
	}

	err := h.walletUsecase.Recharge(c.Context(), req.WalletID, req.Amount)
	if err != nil {
		h.logger.ErrorContext(c.Context(), "failed to recharge wallet", "error", err)
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		sentry.CaptureException(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "recharge successful"})
}

type TransferRequest struct {
	FromWalletID string  `json:"from_wallet_id"`
	ToWalletID   string  `json:"to_wallet_id"`
	Amount       float64 `json:"amount"`
}

func (h *WalletHandler) Transfer(c fiber.Ctx) error {
	var req TransferRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse request"})
	}

	err := h.walletUsecase.Transfer(c.Context(), req.FromWalletID, req.ToWalletID, req.Amount)
	if err != nil {
		h.logger.ErrorContext(c.Context(), "failed to transfer funds", "error", err)
		// Map specific business logic errors to 4xx status codes
		if strings.Contains(err.Error(), "insufficient funds") || strings.Contains(err.Error(), "cannot transfer to the same wallet") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		// Report unexpected errors to Sentry
		sentry.CaptureException(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "transfer successful"})
}
