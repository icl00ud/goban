// Command goban is the unified CLI entry point for the Goban kanban tool.
// It exposes subcommands for the CLI, MCP server, and Telegram bot.
package main

import (
	"os"

	"github.com/icl00ud/goban/pkg/bot"
	"github.com/icl00ud/goban/pkg/cli"
	"github.com/icl00ud/goban/pkg/mcp"
	"github.com/spf13/cobra"
)

func main() {
	// MCP server command (stdio transport)
	mcpCmd := &cobra.Command{
		Use:   "mcp",
		Short: "Start MCP server (stdio transport)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return mcp.Serve()
		},
	}

	// Telegram bot command
	botCmd := &cobra.Command{
		Use:   "bot",
		Short: "Start Telegram bot (long-polling)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return bot.Start()
		},
	}

	cli.RootCmd.AddCommand(mcpCmd, botCmd)

	if err := cli.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
