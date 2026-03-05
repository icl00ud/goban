package client

import "fmt"

// CardResponse mirrors dto.CardResponse.
type CardResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Position    int    `json:"position"`
	Priority    string `json:"priority"`
	ColumnID    uint   `json:"column_id"`
}

// ColumnResponse mirrors dto.ColumnResponse.
type ColumnResponse struct {
	ID       uint           `json:"id"`
	Title    string         `json:"title"`
	Position int            `json:"position"`
	BoardID  uint           `json:"board_id"`
	Cards    []CardResponse `json:"cards,omitempty"`
}

// BoardResponse mirrors dto.BoardResponse.
type BoardResponse struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Color       string           `json:"color"`
	Position    int              `json:"position"`
	Columns     []ColumnResponse `json:"columns,omitempty"`
	CreatedAt   string           `json:"created_at"`
}

// ListBoards returns all boards for the authenticated user.
func (c *Client) ListBoards() ([]BoardResponse, error) {
	return get[[]BoardResponse](c, "/api/v1/boards")
}

// GetBoard returns a single board with its columns and cards.
func (c *Client) GetBoard(id uint) (*BoardResponse, error) {
	return get[*BoardResponse](c, fmt.Sprintf("/api/v1/boards/%d", id))
}

// CreateBoard creates a new board.
func (c *Client) CreateBoard(name, description, color string) (*BoardResponse, error) {
	body := map[string]string{"name": name, "description": description, "color": color}
	return post[*BoardResponse](c, "/api/v1/boards", body)
}

// UpdateBoard updates an existing board.
func (c *Client) UpdateBoard(id uint, name, description, color string) (*BoardResponse, error) {
	body := map[string]string{"name": name, "description": description, "color": color}
	return put[*BoardResponse](c, fmt.Sprintf("/api/v1/boards/%d", id), body)
}

// DeleteBoard deletes a board.
func (c *Client) DeleteBoard(id uint) error {
	_, err := del[any](c, fmt.Sprintf("/api/v1/boards/%d", id))
	return err
}
