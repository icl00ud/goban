package models

import (
	"time"

	"gorm.io/gorm"
)

type Column struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"type:varchar(100);not null" json:"title"`
	Position  int            `gorm:"not null;default:0" json:"position"`
	BoardID   uint           `gorm:"not null;index" json:"board_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Board     Board          `gorm:"foreignKey:BoardID" json:"-"`
	Cards     []Card         `gorm:"foreignKey:ColumnID;constraint:OnDelete:CASCADE" json:"cards,omitempty"`
}
