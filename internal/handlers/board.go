package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/icl00ud/goban/internal/dto"
	"github.com/icl00ud/goban/internal/models"
	"github.com/icl00ud/goban/internal/services"
	"github.com/icl00ud/goban/internal/utils"
)

type BoardHandler struct {
	boardService *services.BoardService
}

func NewBoardHandler(boardService *services.BoardService) *BoardHandler {
	return &BoardHandler{boardService: boardService}
}

// List returns all boards for the authenticated user
func (h *BoardHandler) List(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	boards, err := h.boardService.GetAllByUser(userID)
	if err != nil {
		return utils.InternalError(c, "Failed to fetch boards")
	}

	response := make([]dto.BoardResponse, len(boards))
	for i, board := range boards {
		response[i] = dto.BoardResponse{
			ID:          board.ID,
			Name:        board.Name,
			Description: board.Description,
			Color:       board.Color,
			Position:    board.Position,
			CreatedAt:   board.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return utils.Success(c, response)
}

// Create creates a new board
func (h *BoardHandler) Create(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req dto.CreateBoardRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if req.Name == "" {
		return utils.BadRequest(c, "Board name is required")
	}

	board, err := h.boardService.Create(userID, req.Name, req.Description, req.Color)
	if err != nil {
		return utils.InternalError(c, "Failed to create board")
	}

	return utils.Created(c, toBoardResponse(board))
}

// Get retrieves a single board with its columns and cards
func (h *BoardHandler) Get(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	boardID, err := c.ParamsInt("id")
	if err != nil {
		return utils.BadRequest(c, "Invalid board ID")
	}

	board, err := h.boardService.GetByID(uint(boardID), userID)
	if err != nil {
		if errors.Is(err, services.ErrNotBoardOwner) {
			return utils.Forbidden(c, err.Error())
		}
		if errors.Is(err, services.ErrBoardNotFound) {
			return utils.NotFound(c, err.Error())
		}
		return utils.InternalError(c, "Failed to fetch board")
	}

	return utils.Success(c, toBoardResponse(board))
}

// Update updates a board
func (h *BoardHandler) Update(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	boardID, err := c.ParamsInt("id")
	if err != nil {
		return utils.BadRequest(c, "Invalid board ID")
	}

	var req dto.UpdateBoardRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	board, err := h.boardService.Update(uint(boardID), userID, req.Name, req.Description, req.Color)
	if err != nil {
		if errors.Is(err, services.ErrNotBoardOwner) {
			return utils.Forbidden(c, err.Error())
		}
		if errors.Is(err, services.ErrBoardNotFound) {
			return utils.NotFound(c, err.Error())
		}
		return utils.InternalError(c, "Failed to update board")
	}

	return utils.Success(c, dto.BoardResponse{
		ID:          board.ID,
		Name:        board.Name,
		Description: board.Description,
		Color:       board.Color,
		Position:    board.Position,
		CreatedAt:   board.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// Reorder reorders boards for the authenticated user
func (h *BoardHandler) Reorder(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req dto.ReorderBoardsRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if len(req.BoardIDs) == 0 {
		return utils.BadRequest(c, "board_ids is required")
	}

	err := h.boardService.Reorder(userID, req.BoardIDs)
	if err != nil {
		if errors.Is(err, services.ErrNotBoardOwner) {
			return utils.Forbidden(c, err.Error())
		}
		return utils.InternalError(c, "Failed to reorder boards")
	}

	return utils.SuccessWithMessage(c, "Boards reordered successfully")
}

// Delete deletes a board
func (h *BoardHandler) Delete(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	boardID, err := c.ParamsInt("id")
	if err != nil {
		return utils.BadRequest(c, "Invalid board ID")
	}

	err = h.boardService.Delete(uint(boardID), userID)
	if err != nil {
		if errors.Is(err, services.ErrNotBoardOwner) {
			return utils.Forbidden(c, err.Error())
		}
		return utils.InternalError(c, "Failed to delete board")
	}

	return utils.SuccessWithMessage(c, "Board deleted successfully")
}

// toBoardResponse converts a Board model to BoardResponse DTO
func toBoardResponse(board *models.Board) dto.BoardResponse {
	response := dto.BoardResponse{
		ID:          board.ID,
		Name:        board.Name,
		Description: board.Description,
		Color:       board.Color,
		Position:    board.Position,
		CreatedAt:   board.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if board.Columns != nil {
		response.Columns = make([]dto.ColumnResponse, len(board.Columns))
		for i, col := range board.Columns {
			response.Columns[i] = dto.ColumnResponse{
				ID:       col.ID,
				Title:    col.Title,
				Position: col.Position,
				BoardID:  col.BoardID,
			}

			if col.Cards != nil {
				response.Columns[i].Cards = make([]dto.CardResponse, len(col.Cards))
				for j, card := range col.Cards {
					response.Columns[i].Cards[j] = dto.CardResponse{
						ID:          card.ID,
						Title:       card.Title,
						Description: card.Description,
						Position:    card.Position,
						Priority:    card.Priority,
						ColumnID:    card.ColumnID,
					}
				}
			}
		}
	}

	return response
}
