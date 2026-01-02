package models

import (
	"time"

	"gorm.io/gorm"
)

type Board struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Color       string         `gorm:"type:varchar(7);default:'#3b82f6'" json:"color"`
	Position    int            `gorm:"default:0" json:"position"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	User        User           `gorm:"foreignKey:UserID" json:"-"`
	Columns     []Column       `gorm:"foreignKey:BoardID;constraint:OnDelete:CASCADE" json:"columns,omitempty"`
}
