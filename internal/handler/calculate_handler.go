package handler

import (
	"github.com/gofiber/fiber/v3"

	"github.com/niksis02/pack-calculator-be/internal/model"
	"github.com/niksis02/pack-calculator-be/internal/service"
)

// CalculateHandler handles the pack calculation endpoint.
type CalculateHandler struct {
	svc *service.PackService
}

// NewCalculateHandler creates a CalculateHandler bound to the provided service.
func NewCalculateHandler(svc *service.PackService) *CalculateHandler {
	return &CalculateHandler{svc: svc}
}

// Calculate handles POST /api/v1/calculate.
// Body: {"items": 251}
func (h *CalculateHandler) Calculate(c fiber.Ctx) error {
	var req model.CalculateRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	if req.Items < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "items must be >= 1",
		})
	}
	resp, err := h.svc.Calculate(req.Items)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(resp)
}
