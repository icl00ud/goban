package tools

import (
	"context"
	"fmt"

	"github.com/icl00ud/goban/pkg/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterCardTools registers card-related MCP tools on s.
func RegisterCardTools(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("list_cards",
		mcp.WithDescription("List all cards in a board, grouped by column"),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		boardID := uint(req.GetFloat("board_id", 0))
		board, err := c.GetBoard(boardID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(board.Columns)
	})

	s.AddTool(mcp.NewTool("get_card",
		mcp.WithDescription("Get details of a single card"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("Card ID")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := uint(req.GetFloat("id", 0))
		card, err := c.GetCard(id)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(card)
	})

	s.AddTool(mcp.NewTool("create_card",
		mcp.WithDescription("Create a card in a column"),
		mcp.WithNumber("column_id", mcp.Required(), mcp.Description("Column ID")),
		mcp.WithString("title", mcp.Required(), mcp.Description("Card title")),
		mcp.WithString("description", mcp.Description("Card description")),
		mcp.WithString("priority", mcp.Description("Priority: low | medium | high (default: medium)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		colID := uint(req.GetFloat("column_id", 0))
		title := req.GetString("title", "")
		desc := req.GetString("description", "")
		priority := req.GetString("priority", "medium")
		card, err := c.CreateCard(colID, title, desc, priority)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(card)
	})

	s.AddTool(mcp.NewTool("update_card",
		mcp.WithDescription("Update a card's title, description, or priority"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("Card ID")),
		mcp.WithString("title", mcp.Description("New title")),
		mcp.WithString("description", mcp.Description("New description")),
		mcp.WithString("priority", mcp.Description("New priority: low | medium | high")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := uint(req.GetFloat("id", 0))
		title := req.GetString("title", "")
		desc := req.GetString("description", "")
		priority := req.GetString("priority", "")
		card, err := c.UpdateCard(id, title, desc, priority)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(card)
	})

	s.AddTool(mcp.NewTool("delete_card",
		mcp.WithDescription("Delete a card"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("Card ID")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := uint(req.GetFloat("id", 0))
		if err := c.DeleteCard(id); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Card %d deleted", id)), nil
	})

	s.AddTool(mcp.NewTool("move_card",
		mcp.WithDescription("Move a card to a different column"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("Card ID")),
		mcp.WithNumber("target_column_id", mcp.Required(), mcp.Description("Target column ID")),
		mcp.WithNumber("position", mcp.Description("Position in target column (0-based)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := uint(req.GetFloat("id", 0))
		targetCol := uint(req.GetFloat("target_column_id", 0))
		pos := req.GetInt("position", 0)
		card, err := c.MoveCard(id, targetCol, pos)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(card)
	})
}
