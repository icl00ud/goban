package services

import (
	"errors"

	"github.com/icl00ud/goban/internal/models"
	"github.com/icl00ud/goban/internal/repository"
)

var (
	ErrColumnNotFound = errors.New("column not found")
)

type ColumnService struct {
	columnRepo *repository.ColumnRepository
	boardRepo  *repository.BoardRepository
}

func NewColumnService(columnRepo *repository.ColumnRepository, boardRepo *repository.BoardRepository) *ColumnService {
	return &ColumnService{
		columnRepo: columnRepo,
		boardRepo:  boardRepo,
	}
}

// Create creates a new column at the end of the board
func (s *ColumnService) Create(boardID, userID uint, title string) (*models.Column, error) {
	// Check board ownership
	if !s.boardRepo.BelongsToUser(boardID, userID) {
		return nil, ErrNotBoardOwner
	}

	// Get max position and add to end
	maxPos := s.columnRepo.GetMaxPosition(boardID)

	column := &models.Column{
		Title:    title,
		Position: maxPos + 1,
		BoardID:  boardID,
	}

	if err := s.columnRepo.Create(column); err != nil {
		return nil, err
	}

	return column, nil
}

// GetByID retrieves a column by ID with ownership check
func (s *ColumnService) GetByID(columnID, userID uint) (*models.Column, error) {
	column, err := s.columnRepo.FindByIDWithCards(columnID)
	if err != nil {
		return nil, ErrColumnNotFound
	}

	// Check board ownership
	if !s.boardRepo.BelongsToUser(column.BoardID, userID) {
		return nil, ErrNotBoardOwner
	}

	return column, nil
}

// Update updates a column with ownership check
func (s *ColumnService) Update(columnID, userID uint, title string) (*models.Column, error) {
	column, err := s.columnRepo.FindByID(columnID)
	if err != nil {
		return nil, ErrColumnNotFound
	}

	// Check board ownership
	if !s.boardRepo.BelongsToUser(column.BoardID, userID) {
		return nil, ErrNotBoardOwner
	}

	if title != "" {
		column.Title = title
	}

	if err := s.columnRepo.Update(column); err != nil {
		return nil, err
	}

	return column, nil
}

// Delete deletes a column with ownership check
func (s *ColumnService) Delete(columnID, userID uint) error {
	column, err := s.columnRepo.FindByID(columnID)
	if err != nil {
		return ErrColumnNotFound
	}

	// Check board ownership
	if !s.boardRepo.BelongsToUser(column.BoardID, userID) {
		return ErrNotBoardOwner
	}

	return s.columnRepo.Delete(columnID)
}

// Reorder reorders columns based on the provided order
func (s *ColumnService) Reorder(boardID, userID uint, columnIDs []uint) error {
	// Check board ownership
	if !s.boardRepo.BelongsToUser(boardID, userID) {
		return ErrNotBoardOwner
	}

	// Build positions array
	positions := make([]int, len(columnIDs))
	for i := range columnIDs {
		positions[i] = i
	}

	return s.columnRepo.UpdatePositions(columnIDs, positions)
}

// GetBoardIDForColumn returns the board ID for a column
func (s *ColumnService) GetBoardIDForColumn(columnID uint) (uint, error) {
	column, err := s.columnRepo.FindByID(columnID)
	if err != nil {
		return 0, ErrColumnNotFound
	}
	return column.BoardID, nil
}
