package dto

// CreateColumnRequest represents the request to create a column
type CreateColumnRequest struct {
	Title string `json:"title"`
}

// UpdateColumnRequest represents the request to update a column
type UpdateColumnRequest struct {
	Title string `json:"title"`
}

// ReorderColumnsRequest represents the request to reorder columns
type ReorderColumnsRequest struct {
	BoardID   uint   `json:"board_id"`
	ColumnIDs []uint `json:"column_ids"`
}

// ColumnResponse represents column data in responses
type ColumnResponse struct {
	ID       uint           `json:"id"`
	Title    string         `json:"title"`
	Position int            `json:"position"`
	BoardID  uint           `json:"board_id"`
	Cards    []CardResponse `json:"cards,omitempty"`
}
