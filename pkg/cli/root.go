// Package cli implements the goban command-line interface.
package cli

import (
	"fmt"
	"os"

	"github.com/icl00ud/goban/pkg/client"
	"github.com/icl00ud/goban/pkg/config"
	"github.com/spf13/cobra"
)

var (
	flagServer string
	flagToken  string
	flagJSON   bool

	cfg *config.Config
)

// RootCmd is the top-level cobra command.
var RootCmd = &cobra.Command{
	Use:   "goban",
	Short: "Goban – self-hosted kanban board CLI",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip config loading for completion commands
		if cmd.Name() == "__complete" || cmd.Name() == "help" {
			return nil
		}
		var err error
		cfg, err = config.Load()
		if err != nil {
			return err
		}
		if flagServer != "" {
			cfg.ServerURL = flagServer
		}
		if flagToken != "" {
			cfg.Token = flagToken
		}
		return nil
	},
}

func init() {
	RootCmd.PersistentFlags().StringVar(&flagServer, "server", "", "Goban server URL (overrides config)")
	RootCmd.PersistentFlags().StringVar(&flagToken, "token", "", "JWT token (overrides config)")
	RootCmd.PersistentFlags().BoolVar(&flagJSON, "json", false, "Output as JSON")

	RootCmd.AddCommand(authCmd)
	RootCmd.AddCommand(boardsCmd)
	RootCmd.AddCommand(columnsCmd)
	RootCmd.AddCommand(cardsCmd)
}

// newClient returns an authenticated API client using current config.
func newClient() *client.Client {
	return client.New(cfg.ServerURL, cfg.Token)
}

// requireToken fails fast if no token is configured.
func requireToken(cmd *cobra.Command) error {
	if cfg.Token == "" {
		fmt.Fprintln(os.Stderr, "Not authenticated. Run: goban login")
		os.Exit(1)
	}
	return nil
}
