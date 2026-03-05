package cli

import (
	"fmt"

	"github.com/icl00ud/goban/pkg/cli/format"
	"github.com/spf13/cobra"
)

var boardsCmd = &cobra.Command{
	Use:   "boards",
	Short: "Manage boards",
}

var boardsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all boards",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(cmd); err != nil {
			return err
		}
		boards, err := newClient().ListBoards()
		if err != nil {
			return err
		}
		if flagJSON {
			format.JSON(boards)
			return nil
		}
		rows := make([][]string, len(boards))
		for i, b := range boards {
			rows[i] = []string{fmt.Sprint(b.ID), b.Name, b.Color, b.Description}
		}
		format.Table([]string{"ID", "Name", "Color", "Description"}, rows)
		return nil
	},
}

var boardsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get board details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(cmd); err != nil {
			return err
		}
		var id uint
		fmt.Sscan(args[0], &id)
		board, err := newClient().GetBoard(id)
		if err != nil {
			return err
		}
		if flagJSON {
			format.JSON(board)
			return nil
		}
		format.Msg("Board: %s (ID: %d)", board.Name, board.ID)
		for _, col := range board.Columns {
			format.Msg("  [%d] %s", col.ID, col.Title)
			for _, card := range col.Cards {
				format.Msg("    - [%d] %s [%s]", card.ID, card.Title, card.Priority)
			}
		}
		return nil
	},
}

var (
	boardName  string
	boardDesc  string
	boardColor string
)

var boardsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new board",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(cmd); err != nil {
			return err
		}
		if boardName == "" {
			return fmt.Errorf("--name is required")
		}
		board, err := newClient().CreateBoard(boardName, boardDesc, boardColor)
		if err != nil {
			return err
		}
		if flagJSON {
			format.JSON(board)
			return nil
		}
		format.Msg("Created board #%d: %s", board.ID, board.Name)
		return nil
	},
}

var boardsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a board",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(cmd); err != nil {
			return err
		}
		var id uint
		fmt.Sscan(args[0], &id)
		board, err := newClient().UpdateBoard(id, boardName, boardDesc, boardColor)
		if err != nil {
			return err
		}
		if flagJSON {
			format.JSON(board)
			return nil
		}
		format.Msg("Updated board #%d: %s", board.ID, board.Name)
		return nil
	},
}

var boardsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a board",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(cmd); err != nil {
			return err
		}
		var id uint
		fmt.Sscan(args[0], &id)
		if err := newClient().DeleteBoard(id); err != nil {
			return err
		}
		format.Msg("Deleted board #%d", id)
		return nil
	},
}

func init() {
	boardsCreateCmd.Flags().StringVar(&boardName, "name", "", "Board name (required)")
	boardsCreateCmd.Flags().StringVar(&boardDesc, "description", "", "Board description")
	boardsCreateCmd.Flags().StringVar(&boardColor, "color", "#3b82f6", "Board color (hex)")

	boardsUpdateCmd.Flags().StringVar(&boardName, "name", "", "New board name")
	boardsUpdateCmd.Flags().StringVar(&boardDesc, "description", "", "New description")
	boardsUpdateCmd.Flags().StringVar(&boardColor, "color", "", "New color")

	boardsCmd.AddCommand(boardsListCmd, boardsGetCmd, boardsCreateCmd, boardsUpdateCmd, boardsDeleteCmd)
}
