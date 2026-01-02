package dto

// CreateBoardRequest represents the request to create a board
type CreateBoardRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

// UpdateBoardRequest represents the request to update a board
type UpdateBoardRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

// BoardResponse represents board data in responses
type BoardResponse struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Color       string           `json:"color"`
	Position    int              `json:"position"`
	Columns     []ColumnResponse `json:"columns,omitempty"`
	CreatedAt   string           `json:"created_at"`
}

// ReorderBoardsRequest represents the request to reorder boards
type ReorderBoardsRequest struct {
	BoardIDs []uint `json:"board_ids"`
}
