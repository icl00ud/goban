package repository

import (
	"github.com/icl00ud/goban/internal/models"
	"gorm.io/gorm"
)

type CardRepository struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) *CardRepository {
	return &CardRepository{db: db}
}

// Create creates a new card
func (r *CardRepository) Create(card *models.Card) error {
	return r.db.Create(card).Error
}

// FindByID finds a card by ID
func (r *CardRepository) FindByID(id uint) (*models.Card, error) {
	var card models.Card
	err := r.db.First(&card, id).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}

// FindAllByColumnID finds all cards for a column
func (r *CardRepository) FindAllByColumnID(columnID uint) ([]models.Card, error) {
	var cards []models.Card
	err := r.db.Where("column_id = ?", columnID).Order("position ASC").Find(&cards).Error
	return cards, err
}

// Update updates a card
func (r *CardRepository) Update(card *models.Card) error {
	return r.db.Save(card).Error
}

// Delete soft deletes a card
func (r *CardRepository) Delete(id uint) error {
	return r.db.Delete(&models.Card{}, id).Error
}

// GetMaxPosition returns the maximum position value for cards in a column
func (r *CardRepository) GetMaxPosition(columnID uint) int {
	var maxPos int
	r.db.Model(&models.Card{}).Where("column_id = ?", columnID).Select("COALESCE(MAX(position), -1)").Scan(&maxPos)
	return maxPos
}

// MoveCard moves a card to a new column and position
func (r *CardRepository) MoveCard(cardID, targetColumnID uint, position int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Shift cards at and after the target position in the target column
		if err := tx.Model(&models.Card{}).
			Where("column_id = ? AND position >= ?", targetColumnID, position).
			Update("position", gorm.Expr("position + 1")).Error; err != nil {
			return err
		}

		// Update the card's column and position
		if err := tx.Model(&models.Card{}).
			Where("id = ?", cardID).
			Updates(map[string]interface{}{
				"column_id": targetColumnID,
				"position":  position,
			}).Error; err != nil {
			return err
		}

		return nil
	})
}

// UpdatePositions updates positions for multiple cards in a transaction
func (r *CardRepository) UpdatePositions(cardIDs []uint, positions []int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for i, id := range cardIDs {
			if err := tx.Model(&models.Card{}).Where("id = ?", id).Update("position", positions[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetColumnBoardID returns the board ID for a card's column
func (r *CardRepository) GetColumnBoardID(cardID uint) (uint, error) {
	var card models.Card
	err := r.db.Preload("Column").First(&card, cardID).Error
	if err != nil {
		return 0, err
	}
	return card.Column.BoardID, nil
}
