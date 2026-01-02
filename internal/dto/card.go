package dto

// CreateCardRequest represents the request to create a card
type CreateCardRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

// UpdateCardRequest represents the request to update a card
type UpdateCardRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

// MoveCardRequest represents the request to move a card to a different column
type MoveCardRequest struct {
	TargetColumnID uint `json:"target_column_id"`
	Position       int  `json:"position"`
}

// ReorderCardsRequest represents the request to reorder cards within a column
type ReorderCardsRequest struct {
	ColumnID uint   `json:"column_id"`
	CardIDs  []uint `json:"card_ids"`
}

// CardResponse represents card data in responses
type CardResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Position    int    `json:"position"`
	Priority    string `json:"priority"`
	ColumnID    uint   `json:"column_id"`
}
