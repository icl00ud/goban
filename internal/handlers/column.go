package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/icl00ud/goban/internal/dto"
	"github.com/icl00ud/goban/internal/services"
	"github.com/icl00ud/goban/internal/utils"
)

type ColumnHandler struct {
	columnService *services.ColumnService
}

func NewColumnHandler(columnService *services.ColumnService) *ColumnHandler {
	return &ColumnHandler{columnService: columnService}
}

// Create creates a new column in a board
func (h *ColumnHandler) Create(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	boardID, err := c.ParamsInt("boardId")
	if err != nil {
		return utils.BadRequest(c, "Invalid board ID")
	}

	var req dto.CreateColumnRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if req.Title == "" {
		return utils.BadRequest(c, "Column title is required")
	}

	column, err := h.columnService.Create(uint(boardID), userID, req.Title)
	if err != nil {
		if errors.Is(err, services.ErrNotBoardOwner) {
			return utils.Forbidden(c, err.Error())
		}
		return utils.InternalError(c, "Failed to create column")
	}

	return utils.Created(c, dto.ColumnResponse{
		ID:       column.ID,
		Title:    column.Title,
		Position: column.Position,
		BoardID:  column.BoardID,
	})
}

// Update updates a column
func (h *ColumnHandler) Update(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	columnID, err := c.ParamsInt("id")
	if err != nil {
		return utils.BadRequest(c, "Invalid column ID")
	}

	var req dto.UpdateColumnRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	column, err := h.columnService.Update(uint(columnID), userID, req.Title)
	if err != nil {
		if errors.Is(err, services.ErrNotBoardOwner) {
			return utils.Forbidden(c, err.Error())
		}
		if errors.Is(err, services.ErrColumnNotFound) {
			return utils.NotFound(c, err.Error())
		}
		return utils.InternalError(c, "Failed to update column")
	}

	return utils.Success(c, dto.ColumnResponse{
		ID:       column.ID,
		Title:    column.Title,
		Position: column.Position,
		BoardID:  column.BoardID,
	})
}

// Delete deletes a column
func (h *ColumnHandler) Delete(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	columnID, err := c.ParamsInt("id")
	if err != nil {
		return utils.BadRequest(c, "Invalid column ID")
	}

	err = h.columnService.Delete(uint(columnID), userID)
	if err != nil {
		if errors.Is(err, services.ErrNotBoardOwner) {
			return utils.Forbidden(c, err.Error())
		}
		if errors.Is(err, services.ErrColumnNotFound) {
			return utils.NotFound(c, err.Error())
		}
		return utils.InternalError(c, "Failed to delete column")
	}

	return utils.SuccessWithMessage(c, "Column deleted successfully")
}

// Reorder reorders columns within a board
func (h *ColumnHandler) Reorder(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req dto.ReorderColumnsRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if req.BoardID == 0 || len(req.ColumnIDs) == 0 {
		return utils.BadRequest(c, "Board ID and column IDs are required")
	}

	err := h.columnService.Reorder(req.BoardID, userID, req.ColumnIDs)
	if err != nil {
		if errors.Is(err, services.ErrNotBoardOwner) {
			return utils.Forbidden(c, err.Error())
		}
		return utils.InternalError(c, "Failed to reorder columns")
	}

	return utils.SuccessWithMessage(c, "Columns reordered successfully")
}
