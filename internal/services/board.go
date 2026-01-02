package services

import (
	"errors"

	"github.com/icl00ud/goban/internal/models"
	"github.com/icl00ud/goban/internal/repository"
)

var (
	ErrBoardNotFound  = errors.New("board not found")
	ErrNotBoardOwner  = errors.New("you don't have access to this board")
)

// Default columns for new boards
var defaultColumns = []string{"To Do", "In Progress", "Done"}

type BoardService struct {
	boardRepo  *repository.BoardRepository
	columnRepo *repository.ColumnRepository
}

func NewBoardService(boardRepo *repository.BoardRepository, columnRepo *repository.ColumnRepository) *BoardService {
	return &BoardService{
		boardRepo:  boardRepo,
		columnRepo: columnRepo,
	}
}

// Create creates a new board with default columns
func (s *BoardService) Create(userID uint, name, description, color string) (*models.Board, error) {
	if color == "" {
		color = "#3b82f6"
	}

	// Get next position
	maxPos := s.boardRepo.GetMaxPosition(userID)

	board := &models.Board{
		Name:        name,
		Description: description,
		Color:       color,
		Position:    maxPos + 1,
		UserID:      userID,
	}

	if err := s.boardRepo.Create(board); err != nil {
		return nil, err
	}

	// Create default columns
	columns := make([]models.Column, len(defaultColumns))
	for i, title := range defaultColumns {
		columns[i] = models.Column{
			Title:    title,
			Position: i,
			BoardID:  board.ID,
		}
	}

	if err := s.columnRepo.CreateBatch(columns); err != nil {
		// Board was created but columns failed - log this but don't fail
		// The user can add columns manually
		return board, nil
	}

	// Reload board with columns
	return s.boardRepo.FindByIDWithDetails(board.ID)
}

// GetByID retrieves a board by ID with ownership check
func (s *BoardService) GetByID(boardID, userID uint) (*models.Board, error) {
	if !s.boardRepo.BelongsToUser(boardID, userID) {
		return nil, ErrNotBoardOwner
	}

	board, err := s.boardRepo.FindByIDWithDetails(boardID)
	if err != nil {
		return nil, ErrBoardNotFound
	}

	return board, nil
}

// GetAllByUser retrieves all boards for a user
func (s *BoardService) GetAllByUser(userID uint) ([]models.Board, error) {
	return s.boardRepo.FindAllByUserID(userID)
}

// Update updates a board with ownership check
func (s *BoardService) Update(boardID, userID uint, name, description, color string) (*models.Board, error) {
	if !s.boardRepo.BelongsToUser(boardID, userID) {
		return nil, ErrNotBoardOwner
	}

	board, err := s.boardRepo.FindByID(boardID)
	if err != nil {
		return nil, ErrBoardNotFound
	}

	if name != "" {
		board.Name = name
	}
	board.Description = description
	if color != "" {
		board.Color = color
	}

	if err := s.boardRepo.Update(board); err != nil {
		return nil, err
	}

	return board, nil
}

// Delete deletes a board with ownership check
func (s *BoardService) Delete(boardID, userID uint) error {
	if !s.boardRepo.BelongsToUser(boardID, userID) {
		return ErrNotBoardOwner
	}

	return s.boardRepo.Delete(boardID)
}

// CheckOwnership verifies if a user owns a board
func (s *BoardService) CheckOwnership(boardID, userID uint) bool {
	return s.boardRepo.BelongsToUser(boardID, userID)
}

// Reorder updates the position of boards
func (s *BoardService) Reorder(userID uint, boardIDs []uint) error {
	// Verify all boards belong to the user
	for _, boardID := range boardIDs {
		if !s.boardRepo.BelongsToUser(boardID, userID) {
			return ErrNotBoardOwner
		}
	}

	return s.boardRepo.UpdatePositions(userID, boardIDs)
}
