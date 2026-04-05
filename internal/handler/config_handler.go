package handler

import (
	"github.com/gofiber/fiber/v3"

	"github.com/niksis02/pack-calculator-be/internal/model"
	"github.com/niksis02/pack-calculator-be/internal/service"
)

// ConfigHandler handles pack configuration endpoints.
type ConfigHandler struct {
	svc *service.PackService
}

// NewConfigHandler creates a ConfigHandler bound to the provided service.
func NewConfigHandler(svc *service.PackService) *ConfigHandler {
	return &ConfigHandler{svc: svc}
}

// GetPacks handles GET /api/v1/config/packs.
func (h *ConfigHandler) GetPacks(c fiber.Ctx) error {
	return c.JSON(h.svc.GetPacks())
}

// SetPacks handles POST /api/v1/config/packs.
// Body: {"packs": [250, 500, 1000]}
func (h *ConfigHandler) SetPacks(c fiber.Ctx) error {
	var req model.PackConfig
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	if err := h.svc.SetPacks(req.Packs); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(h.svc.GetPacks())
}
