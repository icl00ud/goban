package cli

import (
	"fmt"

	"github.com/icl00ud/goban/pkg/cli/format"
	"github.com/spf13/cobra"
)

var columnsCmd = &cobra.Command{
	Use:   "columns",
	Short: "Manage columns",
}

var (
	colBoardID uint64
	colTitle   string
)

var columnsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a column in a board",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(cmd); err != nil {
			return err
		}
		if colBoardID == 0 {
			return fmt.Errorf("--board is required")
		}
		if colTitle == "" {
			return fmt.Errorf("--title is required")
		}
		col, err := newClient().CreateColumn(uint(colBoardID), colTitle)
		if err != nil {
			return err
		}
		if flagJSON {
			format.JSON(col)
			return nil
		}
		format.Msg("Created column #%d: %s", col.ID, col.Title)
		return nil
	},
}

var columnsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a column title",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(cmd); err != nil {
			return err
		}
		var id uint
		fmt.Sscan(args[0], &id)
		col, err := newClient().UpdateColumn(id, colTitle)
		if err != nil {
			return err
		}
		if flagJSON {
			format.JSON(col)
			return nil
		}
		format.Msg("Updated column #%d: %s", col.ID, col.Title)
		return nil
	},
}

var columnsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a column",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(cmd); err != nil {
			return err
		}
		var id uint
		fmt.Sscan(args[0], &id)
		if err := newClient().DeleteColumn(id); err != nil {
			return err
		}
		format.Msg("Deleted column #%d", id)
		return nil
	},
}

func init() {
	columnsCreateCmd.Flags().Uint64VarP(&colBoardID, "board", "b", 0, "Board ID (required)")
	columnsCreateCmd.Flags().StringVar(&colTitle, "title", "", "Column title (required)")

	columnsUpdateCmd.Flags().StringVar(&colTitle, "title", "", "New title (required)")

	columnsCmd.AddCommand(columnsCreateCmd, columnsUpdateCmd, columnsDeleteCmd)
}
