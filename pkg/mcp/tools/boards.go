// Package tools registers MCP tool handlers for board operations.
package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/icl00ud/goban/pkg/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterBoardTools registers board-related MCP tools on s.
func RegisterBoardTools(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("list_boards",
		mcp.WithDescription("List all kanban boards"),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		boards, err := c.ListBoards()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(boards)
	})

	s.AddTool(mcp.NewTool("get_board",
		mcp.WithDescription("Get a board with its columns and cards"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("Board ID")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := uint(req.GetFloat("id", 0))
		board, err := c.GetBoard(id)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(board)
	})

	s.AddTool(mcp.NewTool("create_board",
		mcp.WithDescription("Create a new kanban board"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Board name")),
		mcp.WithString("description", mcp.Description("Board description")),
		mcp.WithString("color", mcp.Description("Board color (hex, e.g. #3b82f6)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := req.GetString("name", "")
		desc := req.GetString("description", "")
		color := req.GetString("color", "")
		board, err := c.CreateBoard(name, desc, color)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(board)
	})

	s.AddTool(mcp.NewTool("update_board",
		mcp.WithDescription("Update a board's name, description, or color"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("Board ID")),
		mcp.WithString("name", mcp.Description("New name")),
		mcp.WithString("description", mcp.Description("New description")),
		mcp.WithString("color", mcp.Description("New color")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := uint(req.GetFloat("id", 0))
		name := req.GetString("name", "")
		desc := req.GetString("description", "")
		color := req.GetString("color", "")
		board, err := c.UpdateBoard(id, name, desc, color)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(board)
	})

	s.AddTool(mcp.NewTool("delete_board",
		mcp.WithDescription("Delete a board and all its columns and cards"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("Board ID")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := uint(req.GetFloat("id", 0))
		if err := c.DeleteBoard(id); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Board %d deleted", id)), nil
	})
}

// jsonResult marshals v to a text tool result.
func jsonResult(v any) (*mcp.CallToolResult, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}
