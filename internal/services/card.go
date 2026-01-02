package services

import (
	"errors"

	"github.com/icl00ud/goban/internal/models"
	"github.com/icl00ud/goban/internal/repository"
)

var (
	ErrCardNotFound = errors.New("card not found")
)

type CardService struct {
	cardRepo   *repository.CardRepository
	columnRepo *repository.ColumnRepository
	boardRepo  *repository.BoardRepository
}

func NewCardService(cardRepo *repository.CardRepository, columnRepo *repository.ColumnRepository, boardRepo *repository.BoardRepository) *CardService {
	return &CardService{
		cardRepo:   cardRepo,
		columnRepo: columnRepo,
		boardRepo:  boardRepo,
	}
}

// Create creates a new card at the end of the column
func (s *CardService) Create(columnID, userID uint, title, description, priority string) (*models.Card, error) {
	// Get column to check board ownership
	column, err := s.columnRepo.FindByID(columnID)
	if err != nil {
		return nil, ErrColumnNotFound
	}

	if !s.boardRepo.BelongsToUser(column.BoardID, userID) {
		return nil, ErrNotBoardOwner
	}

	// Validate priority
	if priority == "" {
		priority = models.PriorityMedium
	}
	if !models.ValidatePriority(priority) {
		priority = models.PriorityMedium
	}

	// Get max position and add to end
	maxPos := s.cardRepo.GetMaxPosition(columnID)

	card := &models.Card{
		Title:       title,
		Description: description,
		Priority:    priority,
		Position:    maxPos + 1,
		ColumnID:    columnID,
	}

	if err := s.cardRepo.Create(card); err != nil {
		return nil, err
	}

	return card, nil
}

// GetByID retrieves a card by ID with ownership check
func (s *CardService) GetByID(cardID, userID uint) (*models.Card, error) {
	card, err := s.cardRepo.FindByID(cardID)
	if err != nil {
		return nil, ErrCardNotFound
	}

	// Get column to check board ownership
	column, err := s.columnRepo.FindByID(card.ColumnID)
	if err != nil {
		return nil, ErrColumnNotFound
	}

	if !s.boardRepo.BelongsToUser(column.BoardID, userID) {
		return nil, ErrNotBoardOwner
	}

	return card, nil
}

// Update updates a card with ownership check
func (s *CardService) Update(cardID, userID uint, title, description, priority string) (*models.Card, error) {
	card, err := s.cardRepo.FindByID(cardID)
	if err != nil {
		return nil, ErrCardNotFound
	}

	// Get column to check board ownership
	column, err := s.columnRepo.FindByID(card.ColumnID)
	if err != nil {
		return nil, ErrColumnNotFound
	}

	if !s.boardRepo.BelongsToUser(column.BoardID, userID) {
		return nil, ErrNotBoardOwner
	}

	if title != "" {
		card.Title = title
	}
	card.Description = description
	if priority != "" && models.ValidatePriority(priority) {
		card.Priority = priority
	}

	if err := s.cardRepo.Update(card); err != nil {
		return nil, err
	}

	return card, nil
}

// Delete deletes a card with ownership check
func (s *CardService) Delete(cardID, userID uint) error {
	card, err := s.cardRepo.FindByID(cardID)
	if err != nil {
		return ErrCardNotFound
	}

	// Get column to check board ownership
	column, err := s.columnRepo.FindByID(card.ColumnID)
	if err != nil {
		return ErrColumnNotFound
	}

	if !s.boardRepo.BelongsToUser(column.BoardID, userID) {
		return ErrNotBoardOwner
	}

	return s.cardRepo.Delete(cardID)
}

// Move moves a card to a different column at a specific position
func (s *CardService) Move(cardID, userID, targetColumnID uint, position int) (*models.Card, error) {
	card, err := s.cardRepo.FindByID(cardID)
	if err != nil {
		return nil, ErrCardNotFound
	}

	// Check ownership of source column
	sourceColumn, err := s.columnRepo.FindByID(card.ColumnID)
	if err != nil {
		return nil, ErrColumnNotFound
	}
	if !s.boardRepo.BelongsToUser(sourceColumn.BoardID, userID) {
		return nil, ErrNotBoardOwner
	}

	// Check ownership of target column
	targetColumn, err := s.columnRepo.FindByID(targetColumnID)
	if err != nil {
		return nil, ErrColumnNotFound
	}
	if !s.boardRepo.BelongsToUser(targetColumn.BoardID, userID) {
		return nil, ErrNotBoardOwner
	}

	// Ensure both columns belong to the same board
	if sourceColumn.BoardID != targetColumn.BoardID {
		return nil, errors.New("cannot move card between different boards")
	}

	// Move the card
	if err := s.cardRepo.MoveCard(cardID, targetColumnID, position); err != nil {
		return nil, err
	}

	// Return updated card
	return s.cardRepo.FindByID(cardID)
}

// Reorder reorders cards within a column
func (s *CardService) Reorder(columnID, userID uint, cardIDs []uint) error {
	// Get column to check board ownership
	column, err := s.columnRepo.FindByID(columnID)
	if err != nil {
		return ErrColumnNotFound
	}

	if !s.boardRepo.BelongsToUser(column.BoardID, userID) {
		return ErrNotBoardOwner
	}

	// Build positions array
	positions := make([]int, len(cardIDs))
	for i := range cardIDs {
		positions[i] = i
	}

	return s.cardRepo.UpdatePositions(cardIDs, positions)
}
