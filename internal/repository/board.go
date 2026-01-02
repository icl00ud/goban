package repository

import (
	"github.com/icl00ud/goban/internal/models"
	"gorm.io/gorm"
)

type BoardRepository struct {
	db *gorm.DB
}

func NewBoardRepository(db *gorm.DB) *BoardRepository {
	return &BoardRepository{db: db}
}

// Create creates a new board
func (r *BoardRepository) Create(board *models.Board) error {
	return r.db.Create(board).Error
}

// FindByID finds a board by ID
func (r *BoardRepository) FindByID(id uint) (*models.Board, error) {
	var board models.Board
	err := r.db.First(&board, id).Error
	if err != nil {
		return nil, err
	}
	return &board, nil
}

// FindByIDWithDetails finds a board with its columns and cards
func (r *BoardRepository) FindByIDWithDetails(id uint) (*models.Board, error) {
	var board models.Board
	err := r.db.
		Preload("Columns", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Preload("Columns.Cards", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		First(&board, id).Error
	if err != nil {
		return nil, err
	}
	return &board, nil
}

// FindAllByUserID finds all boards for a user ordered by position
func (r *BoardRepository) FindAllByUserID(userID uint) ([]models.Board, error) {
	var boards []models.Board
	err := r.db.Where("user_id = ?", userID).Order("position ASC, created_at DESC").Find(&boards).Error
	return boards, err
}

// GetMaxPosition returns the highest position for a user's boards
func (r *BoardRepository) GetMaxPosition(userID uint) int {
	var maxPos int
	r.db.Model(&models.Board{}).Where("user_id = ?", userID).Select("COALESCE(MAX(position), -1)").Scan(&maxPos)
	return maxPos
}

// UpdatePositions updates positions for multiple boards in a transaction
func (r *BoardRepository) UpdatePositions(userID uint, boardIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for i, boardID := range boardIDs {
			result := tx.Model(&models.Board{}).
				Where("id = ? AND user_id = ?", boardID, userID).
				Update("position", i)
			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})
}

// Update updates a board
func (r *BoardRepository) Update(board *models.Board) error {
	return r.db.Save(board).Error
}

// Delete soft deletes a board
func (r *BoardRepository) Delete(id uint) error {
	return r.db.Delete(&models.Board{}, id).Error
}

// BelongsToUser checks if a board belongs to a user
func (r *BoardRepository) BelongsToUser(boardID, userID uint) bool {
	var count int64
	r.db.Model(&models.Board{}).Where("id = ? AND user_id = ?", boardID, userID).Count(&count)
	return count > 0
}
