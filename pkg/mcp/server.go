// Package mcp provides the MCP server for Goban.
package mcp

import (
	"fmt"
	"os"

	"github.com/icl00ud/goban/pkg/client"
	"github.com/icl00ud/goban/pkg/config"
	"github.com/icl00ud/goban/pkg/mcp/tools"
	"github.com/mark3labs/mcp-go/server"
)

// Serve starts the MCP server using stdio transport.
// All log output goes to stderr to keep stdout clean for the JSON-RPC protocol.
func Serve() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	if cfg.Token == "" {
		return fmt.Errorf("no token configured; run: goban login")
	}

	c := client.New(cfg.ServerURL, cfg.Token)

	s := server.NewMCPServer(
		"goban",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	tools.RegisterBoardTools(s, c)
	tools.RegisterColumnTools(s, c)
	tools.RegisterCardTools(s, c)

	fmt.Fprintln(os.Stderr, "Goban MCP server starting (stdio)")
	return server.ServeStdio(s)
}
