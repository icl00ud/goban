package models

import (
	"time"

	"gorm.io/gorm"
)

type Card struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"type:varchar(200);not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Position    int            `gorm:"not null;default:0" json:"position"`
	Priority    string         `gorm:"type:varchar(20);default:'medium'" json:"priority"`
	ColumnID    uint           `gorm:"not null;index" json:"column_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Column      Column         `gorm:"foreignKey:ColumnID" json:"-"`
}

// Priority constants
const (
	PriorityLow    = "low"
	PriorityMedium = "medium"
	PriorityHigh   = "high"
)

// ValidatePriority checks if the priority value is valid
func ValidatePriority(priority string) bool {
	switch priority {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return true
	default:
		return false
	}
}
