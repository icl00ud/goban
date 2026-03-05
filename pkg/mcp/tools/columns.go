package tools

import (
	"context"
	"fmt"

	"github.com/icl00ud/goban/pkg/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterColumnTools registers column-related MCP tools on s.
func RegisterColumnTools(s *server.MCPServer, c *client.Client) {
	s.AddTool(mcp.NewTool("create_column",
		mcp.WithDescription("Create a column in a board"),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		mcp.WithString("title", mcp.Required(), mcp.Description("Column title")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		boardID := uint(req.GetFloat("board_id", 0))
		title := req.GetString("title", "")
		col, err := c.CreateColumn(boardID, title)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(col)
	})

	s.AddTool(mcp.NewTool("update_column",
		mcp.WithDescription("Rename a column"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("Column ID")),
		mcp.WithString("title", mcp.Required(), mcp.Description("New title")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := uint(req.GetFloat("id", 0))
		title := req.GetString("title", "")
		col, err := c.UpdateColumn(id, title)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(col)
	})

	s.AddTool(mcp.NewTool("delete_column",
		mcp.WithDescription("Delete a column and all its cards"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("Column ID")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := uint(req.GetFloat("id", 0))
		if err := c.DeleteColumn(id); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Column %d deleted", id)), nil
	})
}
