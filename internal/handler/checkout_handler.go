package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/illusi03/golearn/internal/request"
	"github.com/illusi03/golearn/internal/service"
)

type CheckoutHandler struct {
	transactionService *service.TransactionService
}

func NewCheckoutHandler(transactionService *service.TransactionService) *CheckoutHandler {
	return &CheckoutHandler{
		transactionService: transactionService,
	}
}

func (h *CheckoutHandler) Checkout(c fiber.Ctx) error {
	req := &request.CheckoutRequest{}
	if err := c.Bind().Body(req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Error validation json",
			"error":   err.Error(),
		})
	}

	data, err := h.transactionService.Checkout(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
			"error":   nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Checkout successful",
		"data":    data,
	})
}

func (h *CheckoutHandler) GetAllTransactions(c fiber.Ctx) error {
	data, err := h.transactionService.FindAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
			"error":   nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data fetched successfully",
		"data":    data,
	})
}
