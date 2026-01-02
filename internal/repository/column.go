package repository

import (
	"github.com/icl00ud/goban/internal/models"
	"gorm.io/gorm"
)

type ColumnRepository struct {
	db *gorm.DB
}

func NewColumnRepository(db *gorm.DB) *ColumnRepository {
	return &ColumnRepository{db: db}
}

// Create creates a new column
func (r *ColumnRepository) Create(column *models.Column) error {
	return r.db.Create(column).Error
}

// FindByID finds a column by ID
func (r *ColumnRepository) FindByID(id uint) (*models.Column, error) {
	var column models.Column
	err := r.db.First(&column, id).Error
	if err != nil {
		return nil, err
	}
	return &column, nil
}

// FindByIDWithCards finds a column with its cards
func (r *ColumnRepository) FindByIDWithCards(id uint) (*models.Column, error) {
	var column models.Column
	err := r.db.
		Preload("Cards", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		First(&column, id).Error
	if err != nil {
		return nil, err
	}
	return &column, nil
}

// FindAllByBoardID finds all columns for a board
func (r *ColumnRepository) FindAllByBoardID(boardID uint) ([]models.Column, error) {
	var columns []models.Column
	err := r.db.Where("board_id = ?", boardID).Order("position ASC").Find(&columns).Error
	return columns, err
}

// Update updates a column
func (r *ColumnRepository) Update(column *models.Column) error {
	return r.db.Save(column).Error
}

// Delete soft deletes a column
func (r *ColumnRepository) Delete(id uint) error {
	return r.db.Delete(&models.Column{}, id).Error
}

// GetMaxPosition returns the maximum position value for columns in a board
func (r *ColumnRepository) GetMaxPosition(boardID uint) int {
	var maxPos int
	r.db.Model(&models.Column{}).Where("board_id = ?", boardID).Select("COALESCE(MAX(position), -1)").Scan(&maxPos)
	return maxPos
}

// UpdatePositions updates positions for multiple columns in a transaction
func (r *ColumnRepository) UpdatePositions(columnIDs []uint, positions []int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for i, id := range columnIDs {
			if err := tx.Model(&models.Column{}).Where("id = ?", id).Update("position", positions[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// CreateBatch creates multiple columns at once
func (r *ColumnRepository) CreateBatch(columns []models.Column) error {
	return r.db.Create(&columns).Error
}
