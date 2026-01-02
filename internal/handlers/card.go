package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/icl00ud/goban/internal/dto"
	"github.com/icl00ud/goban/internal/services"
	"github.com/icl00ud/goban/internal/utils"
)

type CardHandler struct {
	cardService *services.CardService
}

func NewCardHandler(cardService *services.CardService) *CardHandler {
	return &CardHandler{cardService: cardService}
}

// Create creates a new card in a column
func (h *CardHandler) Create(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	columnID, err := c.ParamsInt("columnId")
	if err != nil {
		return utils.BadRequest(c, "Invalid column ID")
	}

	var req dto.CreateCardRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if req.Title == "" {
		return utils.BadRequest(c, "Card title is required")
	}

	card, err := h.cardService.Create(uint(columnID), userID, req.Title, req.Description, req.Priority)
	if err != nil {
		if errors.Is(err, services.ErrNotBoardOwner) {
			return utils.Forbidden(c, err.Error())
		}
		if errors.Is(err, services.ErrColumnNotFound) {
			return utils.NotFound(c, err.Error())
		}
		return utils.InternalError(c, "Failed to create card")
	}

	return utils.Created(c, dto.CardResponse{
		ID:          card.ID,
		Title:       card.Title,
		Description: card.Description,
		Position:    card.Position,
		Priority:    card.Priority,
		ColumnID:    card.ColumnID,
	})
}

// Get retrieves a card by ID
func (h *CardHandler) Get(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	cardID, err := c.ParamsInt("id")
	if err != nil {
		return utils.BadRequest(c, "Invalid card ID")
	}

	card, err := h.cardService.GetByID(uint(cardID), userID)
	if err != nil {
		if errors.Is(err, services.ErrNotBoardOwner) {
			return utils.Forbidden(c, err.Error())
		}
		if errors.Is(err, services.ErrCardNotFound) {
			return utils.NotFound(c, err.Error())
		}
		return utils.InternalError(c, "Failed to fetch card")
	}

	return utils.Success(c, dto.CardResponse{
		ID:          card.ID,
		Title:       card.Title,
		Description: card.Description,
		Position:    card.Position,
		Priority:    card.Priority,
		ColumnID:    card.ColumnID,
	})
}

// Update updates a card
func (h *CardHandler) Update(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	cardID, err := c.ParamsInt("id")
	if err != nil {
		return utils.BadRequest(c, "Invalid card ID")
	}

	var req dto.UpdateCardRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	card, err := h.cardService.Update(uint(cardID), userID, req.Title, req.Description, req.Priority)
	if err != nil {
		if errors.Is(err, services.ErrNotBoardOwner) {
			return utils.Forbidden(c, err.Error())
		}
		if errors.Is(err, services.ErrCardNotFound) {
			return utils.NotFound(c, err.Error())
		}
		return utils.InternalError(c, "Failed to update card")
	}

	return utils.Success(c, dto.CardResponse{
		ID:          card.ID,
		Title:       card.Title,
		Description: card.Description,
		Position:    card.Position,
		Priority:    card.Priority,
		ColumnID:    card.ColumnID,
	})
}

// Delete deletes a card
func (h *CardHandler) Delete(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	cardID, err := c.ParamsInt("id")
	if err != nil {
		return utils.BadRequest(c, "Invalid card ID")
	}

	err = h.cardService.Delete(uint(cardID), userID)
	if err != nil {
		if errors.Is(err, services.ErrNotBoardOwner) {
			return utils.Forbidden(c, err.Error())
		}
		if errors.Is(err, services.ErrCardNotFound) {
			return utils.NotFound(c, err.Error())
		}
		return utils.InternalError(c, "Failed to delete card")
	}

	return utils.SuccessWithMessage(c, "Card deleted successfully")
}

// Move moves a card to a different column
func (h *CardHandler) Move(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	cardID, err := c.ParamsInt("id")
	if err != nil {
		return utils.BadRequest(c, "Invalid card ID")
	}

	var req dto.MoveCardRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if req.TargetColumnID == 0 {
		return utils.BadRequest(c, "Target column ID is required")
	}

	card, err := h.cardService.Move(uint(cardID), userID, req.TargetColumnID, req.Position)
	if err != nil {
		if errors.Is(err, services.ErrNotBoardOwner) {
			return utils.Forbidden(c, err.Error())
		}
		if errors.Is(err, services.ErrCardNotFound) || errors.Is(err, services.ErrColumnNotFound) {
			return utils.NotFound(c, err.Error())
		}
		return utils.InternalError(c, "Failed to move card")
	}

	return utils.Success(c, dto.CardResponse{
		ID:          card.ID,
		Title:       card.Title,
		Description: card.Description,
		Position:    card.Position,
		Priority:    card.Priority,
		ColumnID:    card.ColumnID,
	})
}

// Reorder reorders cards within a column
func (h *CardHandler) Reorder(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req dto.ReorderCardsRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if req.ColumnID == 0 || len(req.CardIDs) == 0 {
		return utils.BadRequest(c, "Column ID and card IDs are required")
	}

	err := h.cardService.Reorder(req.ColumnID, userID, req.CardIDs)
	if err != nil {
		if errors.Is(err, services.ErrNotBoardOwner) {
			return utils.Forbidden(c, err.Error())
		}
		if errors.Is(err, services.ErrColumnNotFound) {
			return utils.NotFound(c, err.Error())
		}
		return utils.InternalError(c, "Failed to reorder cards")
	}

	return utils.SuccessWithMessage(c, "Cards reordered successfully")
}
