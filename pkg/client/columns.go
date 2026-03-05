package client

import "fmt"

// CreateColumn creates a new column in a board.
func (c *Client) CreateColumn(boardID uint, title string) (*ColumnResponse, error) {
	body := map[string]string{"title": title}
	return post[*ColumnResponse](c, fmt.Sprintf("/api/v1/boards/%d/columns", boardID), body)
}

// UpdateColumn renames a column.
func (c *Client) UpdateColumn(id uint, title string) (*ColumnResponse, error) {
	body := map[string]string{"title": title}
	return put[*ColumnResponse](c, fmt.Sprintf("/api/v1/columns/%d", id), body)
}

// DeleteColumn deletes a column.
func (c *Client) DeleteColumn(id uint) error {
	_, err := del[any](c, fmt.Sprintf("/api/v1/columns/%d", id))
	return err
}
