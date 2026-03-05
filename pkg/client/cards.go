package client

import "fmt"

// CreateCard creates a new card in a column.
func (c *Client) CreateCard(columnID uint, title, description, priority string) (*CardResponse, error) {
	body := map[string]string{"title": title, "description": description, "priority": priority}
	return post[*CardResponse](c, fmt.Sprintf("/api/v1/columns/%d/cards", columnID), body)
}

// GetCard returns a card by ID.
func (c *Client) GetCard(id uint) (*CardResponse, error) {
	return get[*CardResponse](c, fmt.Sprintf("/api/v1/cards/%d", id))
}

// UpdateCard updates a card.
func (c *Client) UpdateCard(id uint, title, description, priority string) (*CardResponse, error) {
	body := map[string]string{"title": title, "description": description, "priority": priority}
	return put[*CardResponse](c, fmt.Sprintf("/api/v1/cards/%d", id), body)
}

// DeleteCard deletes a card.
func (c *Client) DeleteCard(id uint) error {
	_, err := del[any](c, fmt.Sprintf("/api/v1/cards/%d", id))
	return err
}

// MoveCard moves a card to another column.
func (c *Client) MoveCard(id, targetColumnID uint, position int) (*CardResponse, error) {
	body := map[string]any{"target_column_id": targetColumnID, "position": position}
	return put[*CardResponse](c, fmt.Sprintf("/api/v1/cards/%d/move", id), body)
}

// ListCards returns all cards of a board, grouped per column.
// It uses GetBoard internally since there is no dedicated /cards endpoint.
func (c *Client) ListCards(boardID uint) ([]CardResponse, error) {
	board, err := c.GetBoard(boardID)
	if err != nil {
		return nil, err
	}
	var cards []CardResponse
	for _, col := range board.Columns {
		cards = append(cards, col.Cards...)
	}
	return cards, nil
}
