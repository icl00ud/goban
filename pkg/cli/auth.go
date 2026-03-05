package cli

import (
	"fmt"
	"os"

	"github.com/icl00ud/goban/pkg/cli/format"
	"github.com/icl00ud/goban/pkg/client"
	"github.com/icl00ud/goban/pkg/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var printToken bool

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in and save token",
	RunE: func(cmd *cobra.Command, args []string) error {
		var email string
		fmt.Print("Email: ")
		fmt.Scan(&email)

		fmt.Print("Password: ")
		passBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			return fmt.Errorf("reading password: %w", err)
		}

		c := client.New(cfg.ServerURL, "")
		resp, err := c.Login(email, string(passBytes))
		if err != nil {
			return err
		}

		cfg.Token = resp.Token
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		if printToken {
			fmt.Println(resp.Token)
			return nil
		}

		format.Msg("Logged in as %s (%s)", resp.User.Name, resp.User.Email)
		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Clear saved token",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg.Token = ""
		if err := config.Save(cfg); err != nil {
			return err
		}
		format.Msg("Logged out")
		return nil
	},
}

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current user",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(cmd); err != nil {
			return err
		}
		user, err := newClient().Me()
		if err != nil {
			return err
		}
		if flagJSON {
			format.JSON(user)
			return nil
		}
		format.Table(
			[]string{"ID", "Email", "Name"},
			[][]string{{fmt.Sprint(user.ID), user.Email, user.Name}},
		)
		return nil
	},
}

func init() {
	loginCmd.Flags().BoolVar(&printToken, "print-token", false, "Print the JWT token to stdout")
	authCmd.AddCommand(loginCmd, logoutCmd, whoamiCmd)

	// Top-level shortcuts
	RootCmd.AddCommand(&cobra.Command{
		Use:   "login",
		Short: "Log in (shortcut for auth login)",
		RunE:  loginCmd.RunE,
	})
	RootCmd.AddCommand(&cobra.Command{
		Use:   "logout",
		Short: "Log out (shortcut for auth logout)",
		RunE:  logoutCmd.RunE,
	})
	RootCmd.AddCommand(&cobra.Command{
		Use:   "whoami",
		Short: "Show current user (shortcut for auth whoami)",
		RunE:  whoamiCmd.RunE,
	})
}
